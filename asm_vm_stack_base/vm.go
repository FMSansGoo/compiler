package asm_vm_stack_base

import (
	"fmt"
	"go-compiler/utils"
)

const (
	FrameSize  = 1024
	StackSize  = 1024
	GlobalSize = 65536
)

type VM struct {
	constants []Object
	stack     []Object
	//The stack pointer
	sp           int // Always points to the next value. Top of stack is stack[sp-1]
	globals      []Object
	instructions Instructions

	//  存储栈帧数据
	frames      []*Frame
	framesIndex int
	// 为纯数字函数启用自有指令。默认关闭。
	JitEnabled bool
}

func NewVM(bytecode *Bytecode) *VM {
	return NewVMWithGlobals(bytecode, nil)
}

func NewVMWithGlobals(bytecode *Bytecode, globals []Object) *VM {
	// 把 global 外层当做一个主函数来执行
	globalFn := &CompiledFunctionObject{Instructions: bytecode.Instructions}
	globalClosure := &ClosureObject{Fn: globalFn}
	globalFrame := NewFrame(globalClosure, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = globalFrame

	if globals == nil {
		globals = make([]Object, GlobalSize)
	}

	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]Object, StackSize),
		sp:           0,
		globals:      globals,

		frames:      frames,
		framesIndex: 1,
	}
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	if vm.framesIndex > FrameSize {
		utils.LogError("program make function over size")
	}
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) Run() error {
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip += 1

		ip := vm.currentFrame().ip
		ins := vm.currentFrame().Instructions()

		op := ins[ip]
		opCode := GetOpCodeFromValue(op)
		utils.LogInfo("op opCode vm.sp", op, opCode, vm.sp)
		//debug mode
		for index, object := range vm.stack {
			if object == nil {
				continue
			}
			utils.LogInfo("stack item ", index, object)
		}
		switch opCode {
		case OpCodeConstant:
			constIndex := ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.frameConstants()[constIndex])
			if err != nil {
				return err
			}
		case OpCodePop:
			vm.pop()
		case OpCodeMul, OpCodeAdd, OpCodeSub, OpCodeDiv, OpCodeAddEquals, OpCodeSubEquals, OpCodeMulEquals, OpCodeDivEquals:
			err := vm.executeBinaryOperation(opCode)
			if err != nil {
				return err
			}
		case OpCodeTrue:
			err := vm.push(&BoolObject{Value: true})
			if err != nil {
				return err
			}
		case OpCodeFalse:
			err := vm.push(&BoolObject{Value: false})
			if err != nil {
				return err
			}
		case OpCodeEquals, OpCodeNotEquals, OpCodeGreaterThan, OpCodeGreaterThanEquals, OpCodeLessThan, OpCodeLessThanEquals:
			err := vm.executeComparison(opCode)
			if err != nil {
				return err
			}
		case OpCodeNot:
			operand := vm.pop()

			switch operand {
			case &BoolObject{Value: true}:
				return vm.push(&BoolObject{Value: false})
			case &BoolObject{Value: false}:
				return vm.push(&BoolObject{Value: true})
			case &NullObject{}:
				return vm.push(&BoolObject{Value: true})
			default:
				return vm.push(&BoolObject{Value: false})
			}
		case OpCodeMinus:
			operand := vm.pop()

			numType := NumberObject{}.ValueType()
			if operand.ValueType() != numType {
				return fmt.Errorf("unsupported type for negation: %s", operand.ValueType())
			}

			value := operand.(*NumberObject).Value
			return vm.push(&NumberObject{Value: -value})
		case OpCodeNull:
			err := vm.push(&NullObject{})
			if err != nil {
				return err
			}
		case OpCodeJumpNotTruthy:
			// 这里塞入真实的地址
			pos := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case OpCodeJump:
			pos := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case OpCodeSetGlobal:
			globalIndex := ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()
		case OpCodeGetGlobal:
			globalIndex := ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case OpCodeSetLocal:
			localIndex := ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case OpCodeGetLocal:
			localIndex := ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case OpCodeGetBuiltin:
			builtinIndex := ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			definition := Builtins[builtinIndex]

			err := vm.push(definition.Builtin)
			if err != nil {
				return err
			}
		case OpCodeArray:
			numElements := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case OpCodeDict:
			numElements := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			dict := vm.buildDict(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			err := vm.push(dict)
			if err != nil {
				return err
			}
		case OpCodeObjectCall:
			index := vm.pop()
			arrayDictObject := vm.pop()

			err := vm.executeObjectCallExpression(arrayDictObject, index)
			if err != nil {
				return err
			}
		case OpCodeClosure:
			utils.LogInfo("in OpCodeClosure")
			// 已编译函数在常量池中的索引
			constIndex := int(ReadUint16(ins[ip+1:]))
			// 在栈中等待的自由变量的数量
			numFree := int(ReadUint16(ins[ip+3:]))
			vm.currentFrame().ip += 4
			//
			err := vm.pushFunctionClosure(int(constIndex), int(numFree))
			if err != nil {
				return err
			}
		case OpCodeFunctionCall:
			utils.LogInfo("in OpCodeFunctionCall")
			numArgs := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			err := vm.executeFunctionCall(int(numArgs))
			if err != nil {
				return err
			}

		case OpCodeReturn:
			utils.LogInfo("in OpCodeReturn", vm.sp, vm.stack[vm.sp-1])
			// 这里可以处理一下，return 看看有没有值
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			if frame.resultOverride != nil {
				returnValue = frame.resultOverride
			}

			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case OpCodeGetProperty:
			name := vm.pop()
			obj := vm.pop()
			err := vm.executeGetProperty(obj, name)
			if err != nil {
				return err
			}
		case OpCodeSetProperty:
			value := vm.pop()
			name := vm.pop()
			obj := vm.pop()
			err := vm.executeSetProperty(obj, name, value)
			if err != nil {
				return err
			}
		case OpCodeClass:
			memberCount := int(ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			err := vm.executeClass(memberCount)
			if err != nil {
				return err
			}
		case OpCodeGetFree:
			freeIndex := ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure.Free[freeIndex])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) buildArray(startIndex, endIndex int) Object {
	elements := make([]Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &ArrayObject{Values: elements}
}

func (vm *VM) buildDict(startIndex, endIndex int) Object {
	dictPairs := make(map[DictKeyObject]Object, 0)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		// 只支持 string and number
		switch key.ValueType() {
		case NumberObject{}.ValueType():
			k, ok := key.(*NumberObject)
			if !ok {
				utils.LogError("unusable as dict key: %s", k.ValueType())
			}

			dictPairs[DictKeyObject{Key: NumberObject{Value: k.Value}}] = value
		case StringObject{}.ValueType():
			k, ok := key.(*StringObject)
			if !ok {
				utils.LogError("unusable as dict key: %s", k.ValueType())
			}

			dictPairs[DictKeyObject{Key: StringObject{Value: k.Value}}] = value

		default:
			utils.LogError("unusable as hash key: %s", key.ValueType())
		}

	}

	return &DictObject{Pairs: dictPairs}
}

func (vm *VM) executeFunctionCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]
	switch callee := callee.(type) {
	case *ClosureObject:
		return vm.callFunctionClosure(callee, numArgs)
	// 处理一下内置方法
	case *BuiltinObject:
		return vm.callBuiltin(callee, numArgs)
	case *BoundMethodObject:
		return vm.callBoundMethod(callee, numArgs)
	case *ConstructorObject:
		return vm.callConstructor(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function and non-built-in")
	}
}

func (vm *VM) callBoundMethod(bound *BoundMethodObject, numArgs int) error {
	args := make([]Object, numArgs)
	copy(args, vm.stack[vm.sp-numArgs:vm.sp])
	vm.sp = vm.sp - numArgs - 1

	if err := vm.push(bound.Method); err != nil {
		return err
	}
	if err := vm.push(bound.Receiver); err != nil {
		return err
	}
	for _, arg := range args {
		if err := vm.push(arg); err != nil {
			return err
		}
	}
	return vm.callFunctionClosure(bound.Method, numArgs+1)
}

func (vm *VM) callConstructor(ctor *ConstructorObject, numArgs int) error {
	inst := &InstanceObject{
		Class:  ctor.Class,
		Fields: map[string]Object{},
	}
	vm.applyFieldInits(ctor.Class, inst)

	method, _ := lookupClassMember(ctor.Class, "new")
	args := make([]Object, numArgs)
	copy(args, vm.stack[vm.sp-numArgs:vm.sp])
	vm.sp = vm.sp - numArgs - 1

	closure, ok := method.(*ClosureObject)
	if !ok {
		return vm.push(inst)
	}
	if err := vm.push(closure); err != nil {
		return err
	}
	if err := vm.push(inst); err != nil {
		return err
	}
	for _, arg := range args {
		if err := vm.push(arg); err != nil {
			return err
		}
	}
	if err := vm.callFunctionClosure(closure, numArgs+1); err != nil {
		return err
	}
	vm.currentFrame().resultOverride = inst
	return nil
}

func (vm *VM) applyFieldInits(class *ClassObject, inst *InstanceObject) {
	if class == nil {
		return
	}
	vm.applyFieldInits(class.Super, inst)
	for name, value := range class.FieldInits {
		inst.Fields[name] = value
	}
}

func lookupClassMember(class *ClassObject, name string) (Object, bool) {
	for current := class; current != nil; current = current.Super {
		if current.Members == nil {
			continue
		}
		member, ok := current.Members[name]
		if ok {
			return member, true
		}
	}
	return nil, false
}

func (vm *VM) executeGetProperty(obj Object, nameObj Object) error {
	name, err := objectString(nameObj)
	if err != nil {
		return err
	}
	switch target := obj.(type) {
	case *ClassObject:
		if name == "new" {
			return vm.push(&ConstructorObject{Class: target})
		}
		member, ok := lookupClassMember(target, name)
		if !ok {
			return vm.push(&NullObject{})
		}
		return vm.push(member)
	case *InstanceObject:
		if value, ok := target.Fields[name]; ok {
			return vm.push(value)
		}
		if name == "new" {
			return vm.push(&NullObject{})
		}
		member, ok := lookupClassMember(target.Class, name)
		if !ok {
			return vm.push(&NullObject{})
		}
		if closure, ok := member.(*ClosureObject); ok {
			return vm.push(&BoundMethodObject{Receiver: target, Method: closure})
		}
		return vm.push(member)
	default:
		if obj == nil {
			return fmt.Errorf("property not supported: nil")
		}
		return fmt.Errorf("property not supported: %s", obj.ValueType())
	}
}

func (vm *VM) executeSetProperty(obj Object, nameObj Object, value Object) error {
	name, err := objectString(nameObj)
	if err != nil {
		return err
	}
	switch target := obj.(type) {
	case *InstanceObject:
		if target.Fields == nil {
			target.Fields = map[string]Object{}
		}
		target.Fields[name] = value
		return vm.push(value)
	case *ClassObject:
		if target.Members == nil {
			target.Members = map[string]Object{}
		}
		target.Members[name] = value
		return vm.push(value)
	default:
		if obj == nil {
			return fmt.Errorf("property assignment not supported: nil")
		}
		return fmt.Errorf("property assignment not supported: %s", obj.ValueType())
	}
}

func (vm *VM) executeClass(memberCount int) error {
	members := make([]struct {
		name  string
		value Object
		kind  int
	}, memberCount)
	for i := memberCount - 1; i >= 0; i-- {
		kindObj := vm.pop()
		value := vm.pop()
		nameObj := vm.pop()
		kindNum, ok := kindObj.(*NumberObject)
		if !ok {
			return fmt.Errorf("class member kind is not a number")
		}
		name, err := objectString(nameObj)
		if err != nil {
			return err
		}
		members[i].name = name
		members[i].value = value
		members[i].kind = int(kindNum.Value)
	}

	superObj := vm.pop()
	var super *ClassObject
	switch parent := superObj.(type) {
	case *NullObject:
	case *ClassObject:
		super = parent
	default:
		if superObj != nil {
			return fmt.Errorf("super is not a class: %s", superObj.ValueType())
		}
	}

	nameObj := vm.pop()
	name, err := objectString(nameObj)
	if err != nil {
		return err
	}

	class := &ClassObject{
		Name:       name,
		Super:      super,
		Members:    map[string]Object{},
		FieldInits: map[string]Object{},
	}
	for _, member := range members {
		if member.kind == FieldInitKind {
			class.FieldInits[member.name] = member.value
			continue
		}
		class.Members[member.name] = member.value
	}
	return vm.push(class)
}

func objectString(obj Object) (string, error) {
	switch value := obj.(type) {
	case *StringObject:
		return value.Value, nil
	case StringObject:
		return value.Value, nil
	default:
		if obj == nil {
			return "", fmt.Errorf("expected string, got nil")
		}
		return "", fmt.Errorf("expected string, got %s", obj.ValueType())
	}
}

func (vm *VM) callFunctionClosure(cl *ClosureObject, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		utils.LogError("callFunctionClosure", fmt.Sprintf("wrong number of arguments: want=%d, got=%d", cl.Fn.NumParameters, numArgs))
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	if vm.JitEnabled {
		code := cl.Fn.numericCode()
		if code != nil {
			return vm.callNumeric(cl, code, numArgs)
		}
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (fn *CompiledFunctionObject) numericCode() *NumericCode {
	if fn.numericChecked {
		return fn.Numeric
	}
	fn.numericChecked = true
	fn.Numeric = compileNumeric(fn)
	return fn.Numeric
}

func (vm *VM) callNumeric(cl *ClosureObject, code *NumericCode, numArgs int) error {
	args := make([]float64, numArgs)
	base := vm.sp - numArgs
	for i := 0; i < numArgs; i++ {
		number, ok := vm.stack[base+i].(*NumberObject)
		if !ok {
			frame := NewFrame(cl, vm.sp-numArgs)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + cl.Fn.NumLocals
			return nil
		}
		args[i] = number.Value
	}

	result, err := runNumeric(code, args)
	if err != nil {
		return err
	}
	vm.sp = vm.sp - numArgs - 1
	return vm.push(&NumberObject{Value: result})
}

func (vm *VM) callBuiltin(builtin *BuiltinObject, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Func(args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		utils.LogInfo("callBuiltin", result)
		vm.push(result)
	} else {
		vm.push(&NullObject{})
	}

	return nil
}

func (vm *VM) executeComparison(op OpCode) error {
	right := vm.pop()
	left := vm.pop()

	utils.LogInfo("executeComparison look left right", left, right)
	numType := NumberObject{}.ValueType()

	if left.ValueType() == numType || right.ValueType() == numType {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case OpCodeEquals:
		return vm.push(&BoolObject{Value: right.(*BoolObject).Value == left.(*BoolObject).Value})
	case OpCodeNotEquals:
		return vm.push(&BoolObject{Value: right.(*BoolObject).Value != left.(*BoolObject).Value})
	default:
		return fmt.Errorf("unknown operator: %+v %s %s", op, left.ValueType(), right.ValueType())
	}
}

func (vm *VM) executeIntegerComparison(op OpCode, left, right Object) error {
	leftValue := left.(*NumberObject).Value
	rightValue := right.(*NumberObject).Value

	utils.LogInfo("executeIntegerComparison look left right", left, right)

	switch op {
	case OpCodeEquals:
		return vm.push(&BoolObject{Value: leftValue == rightValue})
	case OpCodeNotEquals:
		return vm.push(&BoolObject{Value: rightValue != leftValue})
	case OpCodeGreaterThan:
		return vm.push(&BoolObject{Value: leftValue > rightValue})
	case OpCodeGreaterThanEquals:
		return vm.push(&BoolObject{Value: leftValue >= rightValue})
	case OpCodeLessThan:
		return vm.push(&BoolObject{Value: leftValue < rightValue})
	case OpCodeLessThanEquals:
		return vm.push(&BoolObject{Value: leftValue <= rightValue})
	default:
		return fmt.Errorf("unknown operator: %+v ", op)
	}
}

func (vm *VM) executeBinaryOperation(op OpCode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.ValueType()
	rightType := right.ValueType()

	numType := NumberObject{}.ValueType()
	strType := StringObject{}.ValueType()

	switch {
	case leftType == numType && rightType == numType:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == strType && rightType == strType:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryAssignmentOperation(op OpCode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.ValueType()
	rightType := right.ValueType()

	numType := NumberObject{}.ValueType()
	strType := StringObject{}.ValueType()

	switch {
	case leftType == numType && rightType == numType:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == strType && rightType == strType:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op OpCode, left, right Object) error {
	leftValue := left.(*NumberObject).Value
	rightValue := right.(*NumberObject).Value

	var result float64

	switch op {
	case OpCodeAdd, OpCodeAddEquals:
		result = leftValue + rightValue
	case OpCodeSub, OpCodeSubEquals:
		result = leftValue - rightValue
	case OpCodeMul, OpCodeMulEquals:
		result = leftValue * rightValue
	case OpCodeDiv, OpCodeDivEquals:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %+v", op)
	}

	return vm.push(&NumberObject{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op OpCode, left, right Object) error {
	if op != OpCodeAdd && op != OpCodeAddEquals {
		return fmt.Errorf("unknown string operator: %+v", op)
	}

	leftValue := left.(*StringObject).Value
	rightValue := right.(*StringObject).Value

	return vm.push(&StringObject{Value: leftValue + rightValue})
}

func (vm *VM) executeObjectCallExpression(left, index Object) error {
	switch {
	case left.ValueType() == ArrayObject{}.ValueType() && index.ValueType() == NumberObject{}.ValueType():
		return vm.executeArrayIndex(left, index)
	case left.ValueType() == DictObject{}.ValueType():
		return vm.executeDictIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.ValueType())
	}
}

func (vm *VM) executeArrayIndex(array, index Object) error {
	arrayObject := array.(*ArrayObject)
	// 这里可能是多个类型
	i := int64(index.(*NumberObject).Value)
	maxNum := int64(len(arrayObject.Values) - 1)

	if i < 0 || i > maxNum {
		obj := &NullObject{}
		return vm.push(obj)
	}

	return vm.push(arrayObject.Values[i])
}

func (vm *VM) executeDictIndex(hash, index Object) error {
	hashObject := hash.(*DictObject)

	var pair Object

	// todo
	// 1.数字、string 、变量
	switch index.ValueType() {
	case NumberObject{}.ValueType():
		key, ok := index.(*NumberObject)
		if !ok {
			return fmt.Errorf("unusable as dict key: %s", index.ValueType())
		}

		pair, ok = hashObject.Pairs[DictKeyObject{Key: NumberObject{Value: key.Value}}]
		if !ok {
			obj := &NullObject{}
			return vm.push(obj)
		}
	case StringObject{}.ValueType():
		key, ok := index.(*StringObject)
		if !ok {
			return fmt.Errorf("unusable as dict key: %s", index.ValueType())
		}

		pair, ok = hashObject.Pairs[DictKeyObject{Key: StringObject{Value: key.Value}}]
		if !ok {
			obj := &NullObject{}
			return vm.push(obj)
		}
	default:
		return fmt.Errorf("unusable as dict key: %s", index.ValueType())
	}

	return vm.push(pair)
}

func (vm *VM) push(o Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) frameConstants() []Object {
	if fn := vm.currentFrame().cl.Fn; fn != nil && fn.Constants != nil {
		return fn.Constants
	}
	return vm.constants
}

func (vm *VM) pushFunctionClosure(constIndex, numFree int) error {
	constant := vm.frameConstants()[constIndex]
	function, ok := constant.(*CompiledFunctionObject)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	free := make([]Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &ClosureObject{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) pop() Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) GetStackTop() Object {
	utils.LogInfo("GetLastStackItem", vm.sp)
	if vm.sp == 0 {
		return vm.GetLastStackItem()
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) GetLastStackItem() Object {
	return vm.stack[vm.sp]
}

func isTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *BoolObject:
		return obj.Value
	case *NullObject:
		return false
	default:
		return true
	}
}

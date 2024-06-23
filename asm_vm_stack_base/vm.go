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
}

func NewVM(bytecode *Bytecode) *VM {

	// 把 global 外层当做一个主函数来执行
	globalFn := &CompiledFunctionObject{Instructions: bytecode.Instructions}
	globalClosure := &ClosureObject{Fn: globalFn}
	globalFrame := NewFrame(globalClosure, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = globalFrame

	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]Object, StackSize),
		sp:           0,
		globals:      make([]Object, GlobalSize),

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
		utils.LogInfo("op opCode", op, opCode)
		// debug mode
		//for _, object := range vm.stack {
		//	if object == nil {
		//		continue
		//	}
		//	utils.LogInfo("stack item", object)
		//}
		switch opCode {
		case OpCodeConstant:
			constIndex := ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case OpCodePop:
			vm.pop()
		case OpCodeMul, OpCodeAdd, OpCodeSub, OpCodeDiv:
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
		//case OpCodeBreak:
		// todo
		// 这里应该直接跳出来循环

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
			//utils.LogInfo("in numArgs", numArgs)
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

			err := vm.push(returnValue)
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
	default:
		return fmt.Errorf("calling non-function and non-built-in")
	}
}

func (vm *VM) callFunctionClosure(cl *ClosureObject, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		utils.LogError("callFunctionClosure", fmt.Sprintf("wrong number of arguments: want=%d, got=%d", cl.Fn.NumParameters, numArgs))
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

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

func (vm *VM) executeBinaryIntegerOperation(op OpCode, left, right Object) error {
	leftValue := left.(*NumberObject).Value
	rightValue := right.(*NumberObject).Value

	var result float64

	switch op {
	case OpCodeAdd:
		result = leftValue + rightValue
	case OpCodeSub:
		result = leftValue - rightValue
	case OpCodeMul:
		result = leftValue * rightValue
	case OpCodeDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %+v", op)
	}

	return vm.push(&NumberObject{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op OpCode, left, right Object) error {
	if op != OpCodeAdd {
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

func (vm *VM) pushFunctionClosure(constIndex, numFree int) error {
	constant := vm.constants[constIndex]
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

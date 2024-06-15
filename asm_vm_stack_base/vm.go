package asm_vm_stack_base

import (
	"fmt"
	"go-compiler/utils"
)

const (
	StackSize  = 1024
	GlobalSize = 65536
)

type VM struct {
	constants []Object
	stack     []Object
	//The stack pointer
	sp           int // Always points to the next value. Top of stack is stack[sp-1]
	globals      []Object
	pa           int // 暂时的计数器寄存器
	instructions Instructions
}

func NewVM(bytecode *Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]Object, StackSize),
		sp:           0,
		globals:      make([]Object, GlobalSize),
		pa:           0,
	}
}

func (vm *VM) Run() error {
	for vm.pa < len(vm.instructions) {
		op := vm.instructions[vm.pa]
		utils.LogInfo("op ", op)
		opCode := GetOpCodeFromValue(op)
		utils.LogInfo("opCode ", opCode)
		switch opCode {
		case OpCodeConstant:
			constIndex := ReadUint16(vm.instructions[vm.pa+1:])
			vm.pa += 2
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
		case OpCodeEquals, OpCodeNotEquals, OpCodeGreaterThan:
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
			pos := int(ReadUint16(vm.instructions[vm.pa+1:]))
			vm.pa += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.pa = pos - 1
			}
		case OpCodeJump:
			pos := int(ReadUint16(vm.instructions[vm.pa+1:]))
			vm.pa = pos - 1
		case OpCodeSetGlobal:
			globalIndex := ReadUint16(vm.instructions[vm.pa+1:])
			vm.pa += 2

			vm.globals[globalIndex] = vm.pop()
		case OpCodeGetGlobal:
			globalIndex := ReadUint16(vm.instructions[vm.pa+1:])
			vm.pa += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case OpCodeArray:
			numElements := int(ReadUint16(vm.instructions[vm.pa+1:]))
			vm.pa += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case OpCodeDict:
			numElements := int(ReadUint16(vm.instructions[vm.pa+1:]))
			vm.pa += 2

			dict := vm.buildDict(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			err := vm.push(dict)
			if err != nil {
				return err
			}
		}
		// 指针 + 1
		vm.pa += 1
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
	dictPairs := make(map[Object]Object, 0)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]
		dictPairs[key] = value
	}

	return &DictObject{Pairs: dictPairs}
}

func (vm *VM) executeComparison(op OpCode) error {
	right := vm.pop()
	left := vm.pop()

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

	switch op {
	case OpCodeEquals:
		return vm.push(&BoolObject{Value: leftValue == rightValue})
	case OpCodeNotEquals:
		return vm.push(&BoolObject{Value: rightValue != leftValue})
	case OpCodeGreaterThan:
		return vm.push(&BoolObject{Value: rightValue > leftValue})
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

func (vm *VM) push(o Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
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

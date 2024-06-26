package asm_vm_stack_base

import "fmt"

type OpCode struct {
	name  string
	value int64
}

var (
	validOpCodes            = []OpCode{}
	OpCodeConstant          = newOpCode("Constant", 1)
	OpCodeAdd               = newOpCode("Add", 2)
	OpCodeSub               = newOpCode("Sub", 3)
	OpCodeMul               = newOpCode("Mul", 4)
	OpCodeDiv               = newOpCode("Div", 5)
	OpCodePop               = newOpCode("Pop", 10)
	OpCodeTrue              = newOpCode("True", 20)
	OpCodeFalse             = newOpCode("False", 21)
	OpCodeNull              = newOpCode("Null", 22)
	OpCodeEquals            = newOpCode("Equals", 23)
	OpCodeNotEquals         = newOpCode("NotEquals", 24)
	OpCodeNot               = newOpCode("Not", 25)
	OpCodeMinus             = newOpCode("Minus", 26)
	OpCodeLessThan          = newOpCode("LessThan", 27)
	OpCodeLessThanEquals    = newOpCode("LessThanEquals", 28)
	OpCodeGreaterThan       = newOpCode("GreaterThan", 29)
	OpCodeGreaterThanEquals = newOpCode("GreateThanEquals", 30)
	OpCodeAddEquals         = newOpCode("AddEquals", 31)
	OpCodeSubEquals         = newOpCode("SubEquals", 32)
	OpCodeMulEquals         = newOpCode("MulEquals", 33)
	OpCodeDivEquals         = newOpCode("DivEquals", 34)

	OpCodeJumpNotTruthy = newOpCode("JumpNotTruthy", 40)
	OpCodeJump          = newOpCode("Jump", 41)
	OpCodeObjectCall    = newOpCode("ObjectCall", 42)
	//OpCodeBreak         = newOpCode("Break", 43)

	OpCodeSetGlobal = newOpCode("SetGlobal", 50)
	OpCodeGetGlobal = newOpCode("GetGlobal", 51)
	OpCodeSetLocal  = newOpCode("SetLocal", 52)
	OpCodeGetLocal  = newOpCode("GetLocal", 53)

	OpCodeArray = newOpCode("Array", 54)
	OpCodeDict  = newOpCode("Dict", 55)

	OpCodeReturn       = newOpCode("Return", 60)
	OpCodeClosure      = newOpCode("Closure", 61)
	OpCodeFunctionCall = newOpCode("FunctionCall", 62)
	// 获取不受作用域约束的变量指令，这玩意儿术语叫自由变量
	OpCodeGetFree = newOpCode("GetFree", 63)
	// 获取内置函数的变量指令
	OpCodeGetBuiltin = newOpCode("GetBuiltin", 64)
)

func newOpCode(name string, value int64) OpCode {
	o := OpCode{name: name, value: value}
	validOpCodes = append(validOpCodes, o)
	return o
}

func OpCodeAll() []OpCode {
	return validOpCodes
}

func (t OpCode) Value() int64 {
	if !t.valid() {
		panic(fmt.Errorf("invalid OpCode: (%+v)", t))
	}
	return t.value
}

func (t OpCode) ValuePtr() *int64 {
	v := t.Value()
	return &v
}

func (t OpCode) Name() string {
	if !t.valid() {
		panic(fmt.Errorf("invalid OpCode: (%+v)", t))
	}
	return t.name
}

func (t OpCode) NamePtr() *string {
	n := t.Name()
	return &n
}

func (t OpCode) valid() bool {
	for _, v := range OpCodeAll() {
		if v == t {
			return true
		}
	}
	return false
}

func GetOpCodeFromName(s string) OpCode {
	for _, v := range OpCodeAll() {
		if v.name == s {
			return v
		}
	}
	panic(fmt.Errorf("invalid Opcode name: (%+v)", s))
}

func GetOpCodeFromValue(op byte) OpCode {
	for _, v := range OpCodeAll() {
		if v.Value() == int64(op) {
			return v
		}
	}
	panic(fmt.Errorf("invalid Opcode value: (%+v)", op))
}

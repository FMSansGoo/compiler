package asm_vm_stack_base

import "fmt"

type OpCode struct {
	name  string
	value int64
}

var (
	validOpCodes = []OpCode{}
	// opCodeType
	OpCodePush     = newOpCode("Push", 1)
	OpCodePop      = newOpCode("Pop", 2)
	OpCodeConstant = newOpCode("Constant", 3)
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

package asm_vm_stack_base

import "fmt"

type Instructions []byte

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpCodeConstant: {OpCodeConstant.Name(), []int{2}},
}

func Lookup(op string) (*Definition, error) {
	opCode := GetOpCodeFromName(op)
	def, ok := definitions[opCode]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

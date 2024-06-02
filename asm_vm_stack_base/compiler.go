package asm_vm_stack_base

type Compiler struct {
	instructions Instructions
	constants    []Object
}

func New() *Compiler {
	return &Compiler{
		instructions: Instructions{},
		constants:    []Object{},
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions Instructions
	Constants    []Object
}

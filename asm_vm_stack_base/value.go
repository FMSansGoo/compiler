package asm_vm_stack_base

type Object interface {
	ValueType() string
}

type NumberObject struct {
	Value float64
}

func (n NumberObject) ValueType() string {
	return "NumberObject"
}

type BoolObject struct {
	Value string
}

func (b BoolObject) ValueType() string {
	return "BoolObject"
}

type StringObject struct {
	Value string
}

func (b StringObject) ValueType() string {
	return "StringObject"
}

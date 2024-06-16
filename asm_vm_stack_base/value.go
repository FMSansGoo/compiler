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
	Value bool
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

type NullObject struct {
}

func (b NullObject) ValueType() string {
	return "NullObject"
}

type ArrayObject struct {
	Values []Object `json:"values"`
}

func (b ArrayObject) ValueType() string {
	return "ArrayObject"
}

type DictObject struct {
	Pairs map[Object]Object `json:"values"`
}

func (b DictObject) ValueType() string {
	return "HashObject"
}

type CompiledFunctionObject struct {
	Instructions  Instructions `json:"instructions"`
	NumLocals     int          `json:"numLocals"`
	NumParameters int          `json:"numParameters"`
}

func (b CompiledFunctionObject) ValueType() string {
	return "CompiledFunctionObject"
}

type ClosureObject struct {
	Fn   *CompiledFunctionObject
	Free []Object
}

func (b ClosureObject) ValueType() string {
	return "ClosureObject"
}

package asm_vm_stack_base

import "fmt"

type Object interface {
	ValueType() string
	Inspect() string
}

type NumberObject struct {
	Value float64
}

func (n NumberObject) ValueType() string {
	return "NumberObject"
}

func (n NumberObject) Inspect() string {
	return fmt.Sprintf("%f", n.Value)
}

type BoolObject struct {
	Value bool
}

func (b BoolObject) ValueType() string {
	return "BoolObject"
}

func (n BoolObject) Inspect() string {
	return fmt.Sprintf("%v", n.Value)
}

type StringObject struct {
	Value string
}

func (n StringObject) Inspect() string {
	return fmt.Sprintf("%s", n.Value)
}

func (b StringObject) ValueType() string {
	return "StringObject"
}

type NullObject struct {
}

func (n NullObject) Inspect() string {
	return "null"
}

func (n NullObject) ValueType() string {
	return "NullObject"
}

type ArrayObject struct {
	Values []Object `json:"values"`
}

func (b ArrayObject) ValueType() string {
	return "ArrayObject"
}

func (b ArrayObject) Inspect() string {
	ns := []string{}
	for i := 0; i < len(b.Values); i++ {
		ns = append(ns, b.Values[i].Inspect())
	}
	return fmt.Sprintf("%+v", ns)
}

type DictKeyObject struct {
	Key Object
}

func (d DictKeyObject) ValueType() string {
	return "DictKeyObject"
}

func (d DictKeyObject) Inspect() string {
	return fmt.Sprintf("%+v", d.Key.Inspect())
}

type DictObject struct {
	Pairs map[DictKeyObject]Object `json:"values"`
}

func (b DictObject) ValueType() string {
	return "HashObject"
}

func (d DictObject) Inspect() string {
	s := "{\n"
	for key, value := range d.Pairs {
		k := fmt.Sprintf("\tkey[%+v]:", key.Inspect())
		v := fmt.Sprintf("value[%+v]\n", value.Inspect())
		s += k
		s += v
	}
	s += "}\n"
	return s
}

type BuiltinFunction func(args ...Object) Object

type BuiltinObject struct {
	Func BuiltinFunction
}

func (d BuiltinObject) Inspect() string {
	return "func"
}

func (b BuiltinObject) ValueType() string {
	return "BuiltinObject"
}

type CompiledFunctionObject struct {
	Instructions  Instructions `json:"instructions"`
	NumLocals     int          `json:"numLocals"`
	NumParameters int          `json:"numParameters"`
}

func (b CompiledFunctionObject) ValueType() string {
	return "CompiledFunctionObject"
}

func (d CompiledFunctionObject) Inspect() string {
	return "compiledFunc"
}

type ClosureObject struct {
	Fn   *CompiledFunctionObject
	Free []Object
}

func (b ClosureObject) ValueType() string {
	return "ClosureObject"
}

func (d ClosureObject) Inspect() string {
	return "ClosureFunc"
}

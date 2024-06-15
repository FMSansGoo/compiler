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

//type DictKvObject struct {
//	Key   Object
//	Value Object
//}
//
//func (b DictKvObject) ValueType() string {
//	return "DictKvObject"
//}

type DictObject struct {
	Pairs map[Object]Object `json:"values"`
}

func (b DictObject) ValueType() string {
	return "HashObject"
}

package main

import "fmt"

// 数值类型枚举
type ValueType struct {
	name  string
	value int64
}

var (
	ValueTypes = []ValueType{}

	ValueTypeError = newValueType("Error", 0)

	// 数值类型
	ValueTypeNumber             = newValueType("Number", 1)
	ValueTypeString             = newValueType("String", 2)
	ValueTypeBoolean            = newValueType("Boolean", 3)
	ValueTypeNull               = newValueType("Null", 4)
	ValueTypeIdentifier         = newValueType("Identifier", 5)
	ValueTypeBinaryExpression   = newValueType("BinaryExpression", 6)
	ValueTypeFunctionExpression = newValueType("FunctionExpression", 7)
	ValueTypeArrayLiteral       = newValueType("ArrayLiteral", 8)
	ValueTypeDictLiteral        = newValueType("DictLiteral", 9)
	ValueTypePropertyAssignment = newValueType("PropertyAssignment", 10)
	ValueTypeUnaryExpression    = newValueType("UnaryExpression", 11)
	ValueTypeMemberExpression   = newValueType("MemberExpression", 12)
	ValueTypeIfStatement        = newValueType("IfStatement", 13)
	ValueTypeWhileStatement     = newValueType("WhileStatement", 14)
	ValueTypeClassExpression    = newValueType("ClassExpression", 15)
	ValueTypeReturnStatement    = newValueType("ReturnStatement", 16)
)

func newValueType(name string, value int64) ValueType {
	o := ValueType{name: name, value: value}
	ValueTypes = append(ValueTypes, o)
	return o
}

func ValueTypeAll() []ValueType {
	return ValueTypes
}

func (t ValueType) valid() bool {
	for _, v := range ValueTypeAll() {
		if v == t {
			return true
		}
	}
	return false
}

func (t ValueType) Value() int64 {
	if !t.valid() {
		panic(fmt.Errorf("invalid ValueType: (%+v)", t))
	}
	return t.value
}

func (t ValueType) ValuePtr() *int64 {
	v := t.Value()
	return &v
}

func (t ValueType) Name() string {
	if !t.valid() {
		panic(fmt.Errorf("invalid ValueType: (%+v)", t))
	}
	return t.name
}

func (t ValueType) NamePtr() *string {
	n := t.Name()
	return &n
}

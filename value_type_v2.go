package main

type AllType interface {
	ValueType() string
}

type Signature struct {
	Name       string  `json:"name"`       // 名字
	ReturnType AllType `json:"returnType"` // 返回类型
	IsStatic   bool    `json:"isStatic"`   // 是否是静态方法
	VarType    string  `json:"varType"`    // const / var
}

type UnKnownType struct {
}

func (s UnKnownType) ValueType() string {
	return "UnKnownType"
}

type StringType struct {
}

func (s StringType) ValueType() string {
	return "StringType"
}

type NumberType struct {
}

func (n NumberType) ValueType() string {
	return "NumberType"
}

type VoidType struct {
}

func (v VoidType) ValueType() string {
	return "VoidType"
}

type NullType struct {
}

func (n NullType) ValueType() string {
	return "NullType"
}

type ArrayType struct {
	ElementType AllType `json:"elementType"`
}

func (a ArrayType) ValueType() string {
	return "ArrayType"
}

type DictType struct {
	KeyType AllType `json:"KeyType"`
	VType   AllType `json:"vType"`
}

func (d DictType) ValueType() string {
	return "DictType"
}

type BooleanType struct {
}

func (b BooleanType) ValueType() string {
	return "BooleanType"
}

type FunctionType struct {
	Params     []Signature `json:"params"`
	ReturnType AllType     `json:"returnType"`
}

func (f FunctionType) ValueType() string {
	return "FunctionType"
}

type ClassType struct {
	MemberSignatures []Signature `json:"memberSignatures"` // 成员方法签名
	SuperType        AllType     `json:"returnType"`
}

func (c ClassType) ValueType() string {
	return "ClassType"
}

type InstanceType struct {
	ClassType ClassType `json:"classType"` //所属类签名
}

func (i InstanceType) ValueType() string {
	return "InstanceType"
}

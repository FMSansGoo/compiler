package main

type AllType interface {
	ValueType() string
}

type Signature struct {
	Name       string  `json:"name"`
	ReturnType AllType `json:"returnType"`
	IsStatic   bool    `json:"isStatic"` // 是否是静态方法
}

type StringType struct {
}

func (s StringType) ValueType() string {
	return "StringType"
}

type NumberType struct {
}

func (s NumberType) ValueType() string {
	return "NumberType"
}

type VoidType struct {
}

func (s VoidType) ValueType() string {
	return "VoidType"
}

type BooleanType struct {
}

func (s BooleanType) ValueType() string {
	return "BooleanType"
}

type FunctionType struct {
	Params     []Signature `json:"params"`
	ReturnType AllType     `json:"returnType"`
}

func (s FunctionType) ValueType() string {
	return "FunctionType"
}

type ClassType struct {
	MemberSignatures []Signature `json:"memberSignatures"` // 成员方法签名
	SuperType        AllType     `json:"returnType"`
}

func (s ClassType) ValueType() string {
	return "ClassType"
}

type InstanceType struct {
	ClassType ClassType `json:"classType"` //所属类签名
}

func (s InstanceType) ValueType() string {
	return "InstanceType"
}

package main

import "fmt"

type AstType struct {
	name  string
	value int64
}

var (
	AstTypes = []AstType{}

	// AstType
	AstTypeProgram        = newAstType("Program", 1)
	AstTypeBlockStatement = newAstType("BlockStatement", 2)
	//VariableDeclaration
	AstTypeVariableDeclaration = newAstType("VariableDeclaration", 3)
	AstTypeNumberLiteral       = newAstType("NumberLiteral", 4)
	AstTypeNullLiteral         = newAstType("NullLiteral", 5)
	AstTypeIdentifier          = newAstType("Identifier", 6)
	AstTypeStringLiteral       = newAstType("StringLiteral", 7)
	AstTypeBooleanLiteral      = newAstType("BooleanLiteral", 8)
	// 赋值
	AstTypeAssignmentExpression = newAstType("AssignmentExpression", 9)
	// BinaryExpression
	AstTypeBinaryExpression = newAstType("BinaryExpression", 10)
	// funciton
	AstTypeFunctionExpression = newAstType("FunctionExpression", 11)
	// array
	AstTypeArrayLiteral = newAstType("ArrayLiteral", 12)
	// dict
	AstTypeDictLiteral = newAstType("DictLiteral", 13)
	// kv
	AstTypePropertyAssignment = newAstType("PropertyAssignment", 14)
	// not unary
	AstTypeUnaryExpression = newAstType("UnaryExpression", 15)
	// dot array
	AstTypeMemberExpression = newAstType("MemberExpression", 16)
	// if
	AstTypeIfStatement = newAstType("IfStatement", 17)
	// while
	AstTypeWhileStatement = newAstType("WhileStatement", 18)
	// for
	AstTypeForStatement = newAstType("ForStatement", 19)
	// ClassExpression
	AstTypeClassExpression = newAstType("ClassExpression", 20)
	// ClassBodyStatement
	AstTypeClassBodyStatement = newAstType("ClassBodyStatement", 21)
	//ClassVariableDeclaration
	AstTypeClassVariableDeclaration = newAstType("ClassVariableDeclaration", 22)
	// return
	AstTypeReturnStatement = newAstType("ReturnStatement", 23)
	// break
	AstTypeBreakStatement = newAstType("BreakStatement", 24)
)

func newAstType(name string, value int64) AstType {
	o := AstType{name: name, value: value}
	AstTypes = append(AstTypes, o)
	return o
}

func AstTypeAll() []AstType {
	return AstTypes
}

func (t AstType) valid() bool {
	for _, v := range AstTypeAll() {
		if v == t {
			return true
		}
	}
	return false
}

func (t AstType) Value() int64 {
	if !t.valid() {
		panic(fmt.Errorf("invalid AstType: (%+v)", t))
	}
	return t.value
}

func (t AstType) ValuePtr() *int64 {
	v := t.Value()
	return &v
}

func (t AstType) Name() string {
	if !t.valid() {
		panic(fmt.Errorf("invalid AstType: (%+v)", t))
	}
	return t.name
}

func (t AstType) NamePtr() *string {
	n := t.Name()
	return &n
}

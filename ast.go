package main

// Node接口
type Node interface {
	Type() string
}

// Program节点结构
type Program struct {
	Body []Node // body属性
}

// 实现Node接口的Type方法
func (p Program) Type() string {
	return "Program"
}

// Identifier节点结构
type Identifier struct {
	Value string // value属性
}

// 实现Node接口的Type方法
func (i Identifier) Type() string {
	return "Identifier"
}

// BooleanLiteral节点结构
type BooleanLiteral struct {
	Value bool // value属性
}

// 实现Node接口的Type方法
func (bl BooleanLiteral) Type() string {
	return "BooleanLiteral"
}

// NumberLiteral节点结构
type NumberLiteral struct {
	Value float64 // value属性
}

// 实现Node接口的Type方法
func (nl NumberLiteral) Type() string {
	return "NumberLiteral"
}

// StringLiteral节点结构
type StringLiteral struct {
	Value string // value属性
}

// 实现Node接口的Type方法
func (sl StringLiteral) Type() string {
	return "StringLiteral"
}

// NullLiteral节点结构
type NullLiteral struct {
}

// 实现Node接口的Type方法
func (nl NullLiteral) Type() string {
	return "NullLiteral"
}

// VariableDeclaration节点结构
type VariableDeclaration struct {
	Kind  string // kind属性
	Name  Node   // name属性
	Value Node   // value属性
}

// 实现Node接口的Type方法
func (vd VariableDeclaration) Type() string {
	return "VariableDeclaration"
}

// AssignmentExpression节点结构
type AssignmentExpression struct {
	Operator string // operator属性
	Left     Node   // left属性
	Right    Node   // right属性
}

// 实现Node接口的Type方法
func (ae AssignmentExpression) Type() string {
	return "AssignmentExpression"
}

// BinaryExpression节点结构
type BinaryExpression struct {
	Operator string // operator属性
	Left     Node   // left属性
	Right    Node   // right属性
}

// 实现Node接口的Type方法
func (be BinaryExpression) Type() string {
	return "BinaryExpression"
}

type BlockStatement struct {
	Body []Node // body属性
}

func (bs BlockStatement) Type() string {
	return "BlockStatement"
}

type FunctionExpression struct {
	Params []Node // params 属性
	Body   Node   // body 属性
}

func (bs FunctionExpression) Type() string {
	return "FunctionExpression"
}

type ReturnStatement struct {
}

func (bs ReturnStatement) Type() string {
	return "ReturnStatement"
}

type ContinueStatement struct {
}

func (bs ContinueStatement) Type() string {
	return "ContinueStatement"
}

type BreakStatement struct {
}

func (bs BreakStatement) Type() string {
	return "BreakStatement"
}

type UnaryExpression struct {
	Operator string // operator属性
	Value    Node   // value 属性
}

func (ue UnaryExpression) Type() string {
	return "UnaryExpression"
}

// 调用的结点
type CallExpression struct {
	Object Node   // object 属性
	Args   []Node // Args 属性
}

func (ce CallExpression) Type() string {
	return "CallExpression"
}

// 数组的节点
type ArrayLiteral struct {
	Values []Node // values 属性
}

func (al ArrayLiteral) Type() string {
	return "ArrayLiteral"
}

type PropertyAssignment struct {
	Key   Node // key
	Value Node // value
}

func (pa PropertyAssignment) Type() string {
	return "PropertyAssignment"
}

type DictLiteral struct {
	Values []Node // key
}

func (dl DictLiteral) Type() string {
	return "DictLiteral"
}

type IfStatement struct {
	Condition  Node
	Consequent Node
	Alternate  Node
}

func (is IfStatement) Type() string {
	return "IfStatement"
}

type ForStatement struct {
	Init   Node
	Test   Node
	Update Node
	Body   Node
}

func (fs ForStatement) Type() string {
	return "ForStatement"
}

type WhileStatement struct {
	Condition Node
	Body      Node
}

func (w WhileStatement) Type() string {
	return "WhileStatement"
}

type ClassBodyStatement struct {
	Body []Node
}

func (cb ClassBodyStatement) Type() string {
	return "ClassBodyStatement"
}

type ClassExpression struct {
	Name       Node
	SuperClass Node
	Body       Node
}

func (cb ClassExpression) Type() string {
	return "ClassExpression"
}

type ClassLiteral struct {
}

func (cb ClassLiteral) Type() string {
	return "ClassLiteral"
}

type MemberExpression struct {
	Object   Node
	Property Node
	// 这里拿来区分 点语法 和 数组语法
	ElementType string
}

func (cb MemberExpression) Type() string {
	return "MemberExpression"
}

package parser

// Node接口
type Node interface {
	Type() string
}

// Program节点结构
type Program struct {
	Body []Node `json:"body"` // body属性
}

// 实现Node接口的Type方法
func (p Program) Type() string {
	return "Program"
}

// Identifier节点结构
type Identifier struct {
	Value string `json:"value"` // value属性
}

// 实现Node接口的Type方法
func (i Identifier) Type() string {
	return "Identifier"
}

// BooleanLiteral节点结构
type BooleanLiteral struct {
	Value bool `json:"value"` // value属性
}

// 实现Node接口的Type方法
func (bl BooleanLiteral) Type() string {
	return "BooleanLiteral"
}

// NumberLiteral节点结构
type NumberLiteral struct {
	Value float64 `json:"value"` // value属性
}

// 实现Node接口的Type方法
func (nl NumberLiteral) Type() string {
	return "NumberLiteral"
}

// StringLiteral节点结构
type StringLiteral struct {
	Value string `json:"value"` // value属性
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
	Kind  string `json:"kind"`  // kind属性
	Name  Node   `json:"name"`  // name属性
	Value Node   `json:"value"` // value属性
}

// 实现Node接口的Type方法
func (vd VariableDeclaration) Type() string {
	return "VariableDeclaration"
}

// ClassDeclaration节点结构
type ClassVariableDeclaration struct {
	Kind  string `json:"kind"`  // kind属性
	Name  Node   `json:"name"`  // name属性
	Value Node   `json:"value"` // value属性
}

// 实现Node接口的Type方法
func (vd ClassVariableDeclaration) Type() string {
	return "ClassVariableDeclaration"
}

// AssignmentExpression节点结构
type AssignmentExpression struct {
	Operator string `json:"operator"` // operator属性
	Left     Node   `json:"left"`     // left属性
	Right    Node   `json:"right"`    // right属性
}

// 实现Node接口的Type方法
func (ae AssignmentExpression) Type() string {
	return "AssignmentExpression"
}

// BinaryExpression节点结构
type BinaryExpression struct {
	Operator string `json:"operator"` // operator属性
	Left     Node   `json:"left"`     // left属性
	Right    Node   `json:"right"`    // right属性
}

// 实现Node接口的Type方法
func (be BinaryExpression) Type() string {
	return "BinaryExpression"
}

type BlockStatement struct {
	Body []Node `json:"body"` // body属性
}

func (bs BlockStatement) Type() string {
	return "BlockStatement"
}

type FunctionExpression struct {
	Params []Node `json:"params"` // params 属性
	Body   Node   `json:"body"`   // body 属性
}

func (bs FunctionExpression) Type() string {
	return "FunctionExpression"
}

type ReturnStatement struct {
	Value Node `json:"value"`
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
	Operator string `json:"operator"` // operator属性
	Value    Node   `json:"value"`    // value 属性
}

func (ue UnaryExpression) Type() string {
	return "UnaryExpression"
}

// 调用的结点
type CallExpression struct {
	Args   []Node `json:"args"`   // Args 属性
	Object Node   `json:"object"` // object 属性
}

func (ce CallExpression) Type() string {
	return "CallExpression"
}

// 数组的节点
type ArrayLiteral struct {
	Values []Node `json:"values"` // values 属性
}

func (al ArrayLiteral) Type() string {
	return "ArrayLiteral"
}

type PropertyAssignment struct {
	Key   Node `json:"key"`   // key
	Value Node `json:"value"` // value
}

func (pa PropertyAssignment) Type() string {
	return "PropertyAssignment"
}

type DictLiteral struct {
	Values []Node `json:"values"` // key
}

func (dl DictLiteral) Type() string {
	return "DictLiteral"
}

type IfStatement struct {
	Condition  Node `json:"condition"`  // condition属性
	Consequent Node `json:"consequent"` // consequent属性
	Alternate  Node `json:"alternate"`  // alternate属性
}

func (is IfStatement) Type() string {
	return "IfStatement"
}

type ForStatement struct {
	Init   Node `json:"init"`   // init属性
	Test   Node `json:"test"`   // test属性
	Update Node `json:"update"` // update属性
	Body   Node `json:"body"`   // body属性
}

func (fs ForStatement) Type() string {
	return "ForStatement"
}

type WhileStatement struct {
	Condition Node `json:"condition"` // condition属性
	Body      Node `json:"body"`      // body属性
}

func (w WhileStatement) Type() string {
	return "WhileStatement"
}

type ClassBodyStatement struct {
	Body []Node `json:"body"` // body属性
}

func (cb ClassBodyStatement) Type() string {
	return "ClassBodyStatement"
}

type ClassExpression struct {
	Name       Node `json:"name"`       // name属性
	SuperClass Node `json:"superClass"` // superClass属性
	Body       Node `json:"body"`       // body属性
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
	Object      Node   `json:"object"`      // object属性
	Property    Node   `json:"property"`    // property属性
	ElementType string `json:"elementType"` // elementType属性，用于区分点语法和数组语法
}

func (cb MemberExpression) Type() string {
	return "MemberExpression"
}

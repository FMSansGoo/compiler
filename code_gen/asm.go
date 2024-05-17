package code_gen

import (
	"fmt"
	"go-compiler/parser"
	"go-compiler/utils"
)

type CodeGenerator struct {
	Asm         string         `json:"asm"`
	Ast         parser.Program `json:"program_ast"`
	Register    *Register      `json:"register"`
	StackOffset int64          `json:"stack_offset"`
	SymbolTable *SymbolTable
}

func NewCodeGenerator(programAst parser.Program) *CodeGenerator {
	c := &CodeGenerator{
		Ast:         programAst,
		SymbolTable: NewSymbolTable(),
		Register:    NewRegister(),
		StackOffset: 3,
	}
	c.InitAsm()
	return c
}

func (this *CodeGenerator) InitAsm() {
	initAsm := ""
	this.Asm += initAsm
	return
}

func (this *CodeGenerator) Visit() {
	if this.Ast.Type() != parser.AstTypeProgram.Name() {
		return
	}
	this.Asm += this.visit(this.Ast)
}

func (this *CodeGenerator) visit(node parser.Node) string {
	var asm string
	switch node.Type() {
	case parser.AstTypeProgram.Name():
		asm = this.visitProgram(node)
	//数组
	case parser.AstTypeArrayLiteral.Name():
		asm = this.visitArrayLiteral(node)
	// 一元表达式
	case parser.AstTypeUnaryExpression.Name():
		asm = this.visitUnaryExpression(node)
	// 二元表达式
	case parser.AstTypeBinaryExpression.Name():
		asm = this.visitBinaryExpression(node)
	// 赋值
	case parser.AstTypeAssignmentExpression.Name():
		asm = this.visitAssignmentExpression(node)
	// 变量定义
	case parser.AstTypeVariableDeclaration.Name():
		asm = this.visitVariableDeclaration(node)
	// 字面量
	case parser.AstTypeNullLiteral.Name():
		fallthrough
	case parser.AstTypeStringLiteral.Name():
		fallthrough
	case parser.AstTypeBooleanLiteral.Name():
		fallthrough
	case parser.AstTypeNumberLiteral.Name():
		asm = this.visitLiteral(node)
	// 变量名
	case parser.AstTypeIdentifier.Name():
		asm = this.visitIdentifier(node)
	default:
		utils.LogError("visitProgram visit item default", node.Type())
		return ""
	}

	return asm
}

func (this *CodeGenerator) visitProgram(node parser.Node) string {
	if node.Type() != parser.AstTypeProgram.Name() {
		return ""
	}

	body := node.(parser.Program).Body
	asm := ""
	for _, item := range body {
		asm += this.visit(item)
	}
	return asm
}

// 表达式
func (this *CodeGenerator) visitBinaryExpression(node parser.Node) string {
	if node.Type() != parser.AstTypeBinaryExpression.Name() {
		return ""
	}
	instructionMap := map[string]string{
		"+":   "add2",
		"-":   "sub2",
		"*":   "mul2",
		"/":   "div2",
		"and": "bool_and",
		"or":  "bool_or",
		">=":  "bool_lget",
		"<=":  "bool_sget",
		"==":  "bool_eq",
		"!=":  "bool_neq",
		">":   "bool_gt",
		"<":   "bool_st",
		"+=":  "num_add_eq",
		"-=":  "num_sub_eq",
		"*=":  "num_mul_eq",
		"/=":  "num_div_eq",
	}
	ins, ok := instructionMap[node.(parser.BinaryExpression).Operator]
	if !ok {
		utils.LogError("invalid visitBinaryExpression", node.(parser.BinaryExpression).Operator)
		return ""
	}
	// left
	left := node.(parser.BinaryExpression).Left
	leftAsm := this.visit(left)

	// right
	right := node.(parser.BinaryExpression).Right
	rightAsm := this.visit(right)
	utils.LogInfo("visitBinaryExpression", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)

	// 退出来寄存器
	rightReg := this.Register.ReturnRegPop()
	leftReg := this.Register.ReturnRegPop()
	// 暂存结果寄存器
	resultReg := this.Register.ReturnRegAlloc()

	return leftAsm + rightAsm + ins + " " + leftReg + " " + rightReg + " " + resultReg + "\n"
}

func (this *CodeGenerator) visitUnaryExpression(node parser.Node) string {
	if node.Type() != parser.AstTypeUnaryExpression.Name() {
		return ""
	}

	// right
	right := node.(parser.UnaryExpression).Value
	rightAsm := this.visit(right)
	utils.LogInfo("visitUnaryExpression", node.(parser.UnaryExpression).Value)

	// 退出来寄存器
	leftReg := this.Register.ReturnRegPop()
	// 暂存结果寄存器
	resultReg := this.Register.ReturnRegAlloc()

	return rightAsm + "bool_not " + " " + leftReg + " " + resultReg + "\n"
}

func (this *CodeGenerator) visitArrayLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeArrayLiteral.Name() {
		return ""
	}
	var asm string
	arrayValues := node.(parser.ArrayLiteral).Values
	var vs []any
	for _, v := range arrayValues {
		vs = append(vs, this.visitArrayLiteralValue(v))
	}
	for _, v := range vs {
		// 设置一个寄存器数值，然后 push 进去
		asm += fmt.Sprintf("set2 %v %v\n", this.Register.ReturnRegAlloc(), v.(string))
		asm += fmt.Sprintf("push %v\n", this.Register.ReturnRegPop())
		this.StackOffset += 2
	}
	// 最后塞一个数组的数据量
	asm += fmt.Sprintf("set2 %v %v\n", this.Register.ReturnRegAlloc(), len(arrayValues))
	return asm
}

func (this *CodeGenerator) visitLiteral(node parser.Node) string {
	value := ""
	switch node.Type() {
	case parser.AstTypeNullLiteral.Name():
		value = "null"
	case parser.AstTypeStringLiteral.Name():
		value = node.(parser.StringLiteral).Value
	case parser.AstTypeNumberLiteral.Name():
		value = fmt.Sprintf("%v", node.(parser.NumberLiteral).Value)
	case parser.AstTypeBooleanLiteral.Name():
		value = fmt.Sprintf("%v", node.(parser.BooleanLiteral).Value)
	default:
		utils.LogError("visitLiteral invalid type", node.Type())
		return ""
	}
	return fmt.Sprintf("set2 %v %v\n", this.Register.ReturnRegAlloc(), value)
}

func (this *CodeGenerator) visitVariableDeclaration(node parser.Node) string {
	if node.Type() != parser.AstTypeVariableDeclaration.Name() {
		return ""
	}

	left := node.(parser.VariableDeclaration).Name
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		this.visit(left)
	default:
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return ""
	}

	// 做一下限制，变量名不为空
	name := node.(parser.VariableDeclaration).Name
	variableName := name.(parser.Identifier).Value
	if len(variableName) == 0 {
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return ""
	}

	right := node.(parser.VariableDeclaration).Value
	rightAsm := this.visit(right)

	// 退出来寄存器
	asm := rightAsm + "push " + this.Register.ReturnRegPop() + "\n"
	return asm
}

func (this *CodeGenerator) visitAssignmentExpression(node parser.Node) string {
	if node.Type() != parser.AstTypeAssignmentExpression.Name() {
		return ""
	}
	var leftAsm string
	left := node.(parser.AssignmentExpression).Left
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		leftAsm = this.visit(left)
	default:
		utils.LogError("visitAssignmentExpression invalid left variable declaration", left)
		return ""
	}

	// 做一下限制，变量名不为空
	name := node.(parser.AssignmentExpression).Left
	variableName := name.(parser.Identifier).Value
	if len(variableName) == 0 {
		utils.LogError("visitAssignmentExpression invalid left variable declaration", left)
		return ""
	}
	variable, ok := this.SymbolTable.LookupVariableInfo(variableName)
	if !ok {
		utils.LogError("visitAssignmentExpression invalid left variableName", variable)
		return ""
	}
	var asm string
	// 拿到变量的栈指针，设置数值到变量的栈上
	if variable.Address > 0 {
		asm += "set2 f1 " + fmt.Sprintf("%v", variable.Address) + "\n"
	}
	right := node.(parser.AssignmentExpression).Right
	rightAsm := this.visit(right)
	asm = leftAsm + rightAsm + asm + "save_from_register2 " + this.Register.ReturnRegPop() + " f1\n"
	return asm
}

func (this *CodeGenerator) visitIdentifier(node parser.Node) string {
	if node.Type() != parser.AstTypeIdentifier.Name() {
		return ""
	}
	variableName := node.(parser.Identifier).Value
	varInfo, ok := this.SymbolTable.LookupVariableInfo(variableName)
	if !ok {
		this.StackOffset += 2
		this.SymbolTable.AddVariableInfo(variableName, this.StackOffset, true)
		return ""
	}

	var asm string
	// 将数据从栈上拿出来
	offset := varInfo.Address
	asm = "set2 f1 " + fmt.Sprintf("%v", offset) + "\n"
	asm += "load_from_register2 " + "f1 " + this.Register.ReturnRegAlloc() + "\n"
	return asm
}

func (this *CodeGenerator) visitArrayLiteralValue(node parser.Node) string {
	value := ""
	switch node.Type() {
	case parser.AstTypeNullLiteral.Name():
		value = "null"
	case parser.AstTypeStringLiteral.Name():
		value = node.(parser.StringLiteral).Value
	case parser.AstTypeNumberLiteral.Name():
		value = fmt.Sprintf("%v", node.(parser.NumberLiteral).Value)
	case parser.AstTypeBooleanLiteral.Name():
		value = fmt.Sprintf("%v", node.(parser.BooleanLiteral).Value)
	default:
		utils.LogError("visitLiteral invalid type", node.Type())
		return ""
	}
	return value
}

func (this *CodeGenerator) visitNullLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeNullLiteral.Name() {
		return ""
	}

	return "null"
}

func (this *CodeGenerator) visitNumberLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeNumberLiteral.Name() {
		return ""
	}
	value := node.(parser.NumberLiteral).Value
	return fmt.Sprintf("%v", value)
}

func (this *CodeGenerator) visitStringLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeStringLiteral.Name() {
		return ""
	}
	value := node.(parser.StringLiteral).Value
	return fmt.Sprintf("%v", value)
}

func (this *CodeGenerator) visitBooleanLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeBooleanLiteral.Name() {
		return ""
	}
	value := node.(parser.BooleanLiteral).Value
	return fmt.Sprintf("%v", value)
}

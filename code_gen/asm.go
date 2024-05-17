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
	// 表达式
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
	// op
	switch node.(parser.BinaryExpression).Operator {
	case "+":
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

		return leftAsm + rightAsm + "add2 " + leftReg + " " + rightReg + " " + resultReg + "\n"
	case "-":
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

		return leftAsm + rightAsm + "sub2 " + leftReg + " " + rightReg + " " + resultReg + "\n"
	case "*":
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

		return leftAsm + rightAsm + "mul2 " + leftReg + " " + rightReg + " " + resultReg + "\n"
	case "/":
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

		return leftAsm + rightAsm + "div2 " + leftReg + " " + rightReg + " " + resultReg + "\n"
	case "and":
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

		return leftAsm + rightAsm + "bool_and " + leftReg + " " + rightReg + " " + resultReg + "\n"
	case "or":
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

		return leftAsm + rightAsm + "bool_or " + leftReg + " " + rightReg + " " + resultReg + "\n"
	default:
		utils.LogError("visitBinaryExpression invalid operator", node.(parser.BinaryExpression).Operator)
	}
	return ""
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
	asm := rightAsm + "push2 " + this.Register.ReturnRegPop() + "\n"
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

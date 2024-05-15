package code_gen

import (
	"fmt"
	"go-compiler/parser"
	"go-compiler/utils"
)

type CodeGenerator struct {
	Asm         string         `json:"asm"`
	Ast         parser.Program `json:"program_ast"`
	StackOffset int64          `json:"stack_offset"`
	SymbolTable *SymbolTable
}

func NewCodeGenerator(programAst parser.Program) *CodeGenerator {
	c := &CodeGenerator{
		Ast:         programAst,
		Asm:         "",
		SymbolTable: NewSymbolTable(),
		StackOffset: 3,
	}
	c.InitAsm()
	return c
}

func (this *CodeGenerator) InitAsm() {
	initAsm := "; 栈是前1024字节\njump @1024\n.memory 1024\nset2 f1 3\n"
	this.Asm += initAsm
	return
}

func (this *CodeGenerator) Visit() {
	if this.Ast.Type() != parser.AstTypeProgram.Name() {
		return
	}
	this.visitProgram(this.Ast.Body)
}

func (this *CodeGenerator) visitProgram(body []parser.Node) {
	for _, item := range body {
		utils.LogInfo("visitProgram visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case parser.AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		default:
			utils.LogError("visitProgram visit item default", item.Type())
		}
	}
}

func (this *CodeGenerator) visitVariableDeclaration(node parser.Node) {
	//var a = 1
	//var b = 2
	//这种就是栈上
	//搞个表，记录变量对应的栈上位置
	//然后把栈上变量读到寄存器上，再操作寄存器

	//比如
	//a 在 stack[0]
	//b在 stack[1]
	//
	//那 a+b 可以编译成
	//load_from_stack @0 a1
	//load_from_stack @1 a2
	//add a1 a2 a3

	//赋值也是一样的

	left := node.(parser.VariableDeclaration).Name
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		this.visitIdentifier(left)
	default:
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}

	// 做一下限制，变量名不为空
	name := node.(parser.VariableDeclaration).Name
	variableName := name.(parser.Identifier).Value
	if len(variableName) == 0 {
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}

	right := node.(parser.VariableDeclaration).Value

	var subAsm string
	// 初始化 asm
	subAsm = "set2 a3 0\n"
	var numberValue string
	switch right.Type() {
	// 访问表达式
	case parser.AstTypeBinaryExpression.Name():
		asm := this.visitBinaryExpression(right)
		subAsm += this.setVariableAsm(asm, "register")
	// 访问函数，函数的返回类型
	case parser.AstTypeNumberLiteral.Name():
		numberValue = this.visitNumberLiteral(right)
		subAsm += this.setVariableAsm(numberValue, "number")
	default:
		utils.LogError("visitVariableDeclaration invalid right variable declaration", right)
		return
	}

	this.StackOffset += 2
	// 这里的 address 不能写死
	this.SymbolTable.AddVariableInfo(variableName, this.StackOffset, true)
	this.Asm += subAsm
	return
}

func (this *CodeGenerator) setVariableAsm(value string, valueType string) (code string) {
	//; 在 f1 所在地址存放值，存放后移动 f1
	//set2 a1 {}
	//save_from_register2 a1 f1
	//; f1 += 2
	//set2 a3 2
	//add2 f1 a3 f1
	code = fmt.Sprintf("; 在 f1 所在地址存放值，存放后移动 f1\n")
	if valueType == "number" {
		code += fmt.Sprintf("set2 a3 %v\n", value)
	} else if valueType == "register" {
		code += value
	}
	code += fmt.Sprintf("save_from_register2 a3 f1\n")
	code += fmt.Sprintf("; f1 += 2\n")
	code += fmt.Sprintf("set2 a1 2\n")
	code += fmt.Sprintf("add2 f1 a1 f1\n")
	return
}

// 表达式
func (this *CodeGenerator) visitBinaryExpression(node parser.Node) string {
	// BinaryExpression节点结构
	//type BinaryExpression struct {
	//	Operator string // operator属性
	//	Left     Node   // left属性
	//	Right    Node   // right属性
	//}
	if node.Type() != parser.AstTypeBinaryExpression.Name() {
		return ""
	}
	// op
	switch node.(parser.BinaryExpression).Operator {
	case "+":
		//set a1 1
		//set a2 2
		//add a1 a2 a3
		//
		//set a1 a3
		//set a2 3
		//add a1 a2 a3

		asm := ""
		leftAsm := ""
		// 这里接受多个类型
		left := node.(parser.BinaryExpression).Left
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftAsm += this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			num := this.visitNumberLiteral(left)
			leftAsm += fmt.Sprintf("set2 a1 %v\n", num)
			leftAsm += "push a1\n"
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return ""
		}

		// right
		right := node.(parser.BinaryExpression).Right
		rightAsm := ""
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightAsm += this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			num := this.visitNumberLiteral(right)
			rightAsm += fmt.Sprintf("set2 a1 %v\n", num)
			rightAsm += "push a1\n"
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return ""
		}
		//a = 1 + 1
		utils.LogInfo("visitBinaryExpression", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)
		asm += fmt.Sprintf("%v%v\n", leftAsm, rightAsm)
		asm += fmt.Sprintf("pop a2\n")
		asm += fmt.Sprintf("pop a1\n")
		asm += fmt.Sprintf("add a1 a2 a3\n")
		asm += fmt.Sprintf("push a3\n")
		return asm
	case "*":

		asm := ""
		leftAsm := ""
		// 这里接受多个类型
		left := node.(parser.BinaryExpression).Left
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftAsm += this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			num := this.visitNumberLiteral(left)
			leftAsm += fmt.Sprintf("set2 a1 %v\n", num)
			leftAsm += "push a1\n"
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return ""
		}

		// right
		right := node.(parser.BinaryExpression).Right
		rightAsm := ""
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightAsm += this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			num := this.visitNumberLiteral(right)
			rightAsm += fmt.Sprintf("set2 a1 %v\n", num)
			rightAsm += "push a1\n"
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return ""
		}
		//a = 1 * 1
		utils.LogInfo("visitBinaryExpression", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)
		asm += fmt.Sprintf("%v%v\n", leftAsm, rightAsm)
		asm += fmt.Sprintf("pop a2\n")
		asm += fmt.Sprintf("pop a1\n")
		asm += fmt.Sprintf("multiply2 a1 a2 a3\n")
		asm += fmt.Sprintf("push a3\n")
		return asm
	default:
		utils.LogError("visitBinaryExpression invalid operator", node.(parser.BinaryExpression).Operator)
	}
	return ""
}

func (this *CodeGenerator) visitNullLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeNullLiteral.Name() {
		return ""
	}
	return fmt.Sprintf("%v", "null")
}

func (this *CodeGenerator) visitStringLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeStringLiteral.Name() {
		return ""
	}
	return fmt.Sprintf("%v", node.(parser.StringLiteral).Value)
}

func (this *CodeGenerator) visitBooleanLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeBooleanLiteral.Name() {
		return ""
	}
	return fmt.Sprintf("%v", node.(parser.BooleanLiteral).Value)
}

func (this *CodeGenerator) visitNumberLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeNumberLiteral.Name() {
		return ""
	}
	utils.LogInfo("visitNumberLiteral", node.(parser.NumberLiteral))
	return fmt.Sprintf("%v", node.(parser.NumberLiteral).Value)
}

func (this *CodeGenerator) visitIdentifier(node parser.Node) (variableInfo VariableInfo, ok bool) {
	if node.Type() != parser.AstTypeIdentifier.Name() {
		return VariableInfo{}, false
	}

	variableInfo, ok = this.SymbolTable.LookupVariableInfo(node.(parser.Identifier).Value)
	if ok {
		return variableInfo, false
	}
	return VariableInfo{
		Name: node.(parser.Identifier).Value,
	}, false
}

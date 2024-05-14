package code_gen

import (
	"go-compiler/parser"
	"go-compiler/utils"
)

type CodeGenerator struct {
	Asm string         `json:"asm"`
	Ast parser.Program `json:"program_ast"`
}

func NewCodeGenerator(programAst parser.Program) *CodeGenerator {
	asm := "jump @1024\n.memory 1024\nset2 f1 3"
	return &CodeGenerator{
		Ast: programAst,
		Asm: asm,
	}
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
	////type VariableDeclaration struct {
	////	Kind  string // kind属性
	////	Name  Node   // name属性
	////	Value Node   // value属性
	////}
	//left := node.(parser.VariableDeclaration).Name
	//var variableName string
	//switch left.Type() {
	//case parser.AstTypeIdentifier.Name():
	//	_, variableName, _ = this.visitIdentifier(left)
	//default:
	//	utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
	//	return
	//}
	return
}

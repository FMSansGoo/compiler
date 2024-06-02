package asm_vm_register_base

import (
	"fmt"
	"go-compiler/parser"
	"go-compiler/utils"
)

type FunctionInfo struct {
	Name      string `json:"name"`
	RetOffset int64  `json:"ret_offset"`
}

type CodeGenerator struct {
	Asm         string         `json:"asm"`
	Ast         parser.Program `json:"program_ast"`
	Register    *Register      `json:"register"`
	StackOffset int64          `json:"stack_offset"`
	SymbolTable *SymbolTable
	// 计数器简单点做
	IfCounter    int64 `json:"if_counter"`
	ElseCounter  int64 `json:"else_counter"`
	WhileCounter int64 `json:"while_counter"`
	ForCounter   int64 `json:"for_counter"`
	// function
	FunctionInfo   map[string]*FunctionInfo `json:"function_info"`
	FunctionOffset int64                    `json:"function_offset"`
}

func NewCodeGenerator(programAst parser.Program) *CodeGenerator {
	c := &CodeGenerator{
		Ast:            programAst,
		SymbolTable:    NewSymbolTable(),
		Register:       NewRegister(),
		StackOffset:    3,
		IfCounter:      1,
		ElseCounter:    1,
		WhileCounter:   1,
		ForCounter:     1,
		FunctionInfo:   make(map[string]*FunctionInfo, 0),
		FunctionOffset: 2,
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
	// call func
	case parser.AstTypeCallExpression.Name():
		asm = this.visitCallExpression(node)
	// return
	case parser.AstTypeReturnStatement.Name():
		asm = this.visitReturnStatement(node)
	// block
	case parser.AstTypeBlockStatement.Name():
		asm = this.visitBlockStatement(node)
	//while
	case parser.AstTypeWhileStatement.Name():
		asm = this.visitWhileStatement(node)
	//for
	case parser.AstTypeForStatement.Name():
		asm = this.visitForStatement(node)
	// if
	case parser.AstTypeIfStatement.Name():
		asm = this.visitIfStatement(node)
	//数组
	case parser.AstTypeArrayLiteral.Name():
		asm = this.visitArrayLiteral(node)
	// function
	case parser.AstTypeFunctionExpression.Name():
		asm = this.visitFunctionExpression(node)
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

func (this *CodeGenerator) visitCallExpression(node parser.Node) string {
	if node.Type() != parser.AstTypeCallExpression.Name() {
		return ""
	}

	var asm string

	object := node.(parser.CallExpression).Object
	var funcName string
	if object.Type() == parser.AstTypeIdentifier.Name() {
		funcName = this.visitRawLiteralValue(object)
		asm += fmt.Sprintf(".call @%s ", funcName)
	}

	args := node.(parser.CallExpression).Args
	argsValue := ""
	for _, arg := range args {
		v := this.visitRawLiteralValue(arg)
		argsValue += v + " "
	}
	asm += argsValue + "\n"

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
		"+":   InstructionAdd2.Name(),
		"-":   InstructionSubtract2.Name(),
		"*":   InstructionMultiply2.Name(),
		"/":   InstructionDiv2.Name(),
		"and": InstructionBoolAnd.Name(),
		"or":  InstructionBoolOr.Name(),
		">=":  InstructionBoolGreaterThanEquals.Name(),
		"<=":  InstructionBoolLessThanEquals.Name(),
		"==":  InstructionBoolEquals.Name(),
		"!=":  InstructionBoolNotEquals.Name(),
		">":   InstructionBoolGreaterThan.Name(),
		"<":   InstructionBoolLessThan.Name(),
		"+=":  InstructionPlusAssign.Name(),
		"-=":  InstructionSubtractAssign.Name(),
		"*=":  InstructionMultiplyAssign.Name(),
		"/=":  InstructionDivideAssign.Name(),
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

func (this *CodeGenerator) visitBlockStatement(node parser.Node) string {
	if node.Type() != parser.AstTypeBlockStatement.Name() {
		return ""
	}
	body := node.(parser.BlockStatement).Body

	var asm string
	for _, item := range body {
		asm += this.visit(item)
	}
	return asm
}

func (this *CodeGenerator) visitReturnStatement(node parser.Node) string {
	if node.Type() != parser.AstTypeReturnStatement.Name() {
		return ""
	}
	value := node.(parser.ReturnStatement).Value

	var asm string
	rawValue := this.visitRawLiteralValue(value)
	asm += fmt.Sprintf(".return %v\n", rawValue)
	return asm
}

func (this *CodeGenerator) visitIfStatement(node parser.Node) string {
	//if a1 @xxx地址
	// 其实你可以直接判断是否为 true
	// 或者是否为 true 值

	if node.Type() != parser.AstTypeIfStatement.Name() {
		return ""
	}

	// 先访问 condition
	condition := node.(parser.IfStatement).Condition
	ifPreAsm := this.visit(condition)

	var elseFlag bool
	// 检测一下 else 是否为空
	alternate := node.(parser.IfStatement).Alternate
	if alternate != nil {
		elseFlag = true
	}

	ifCounter := this.IfCounter
	elseCounter := this.ElseCounter
	ifPreAsm += fmt.Sprintf("if %v @if_block_%v\n", this.Register.ReturnRegPop(), ifCounter)
	if elseFlag {
		ifPreAsm += fmt.Sprintf("jump @else_block_%v\n", elseCounter)
	}
	ifPreAsm += fmt.Sprintf("@if_block_%v\n", ifCounter)

	this.IfCounter += 1
	consequent := node.(parser.IfStatement).Consequent
	blockAsm := this.visit(consequent)
	blockAsm += fmt.Sprintf("jump @if_block_end_%v\n", ifCounter)

	var elseAsm string
	if elseFlag {
		elseAsm += fmt.Sprintf("@else_block_%v\n", elseCounter)
		this.ElseCounter += 1
		elseAsm += this.visit(alternate)
	}

	var asm string
	asm += ifPreAsm
	asm += blockAsm
	asm += elseAsm
	asm += fmt.Sprintf("@if_block_end_%v\n", ifCounter)

	return asm
}

func (this *CodeGenerator) visitWhileStatement(node parser.Node) string {
	if node.Type() != parser.AstTypeWhileStatement.Name() {
		return ""
	}

	// 先访问 condition
	condition := node.(parser.WhileStatement).Condition

	whileCounter := this.WhileCounter
	var whilePreAsm string
	whilePreAsm += fmt.Sprintf("@while_init_%v\n", whileCounter)
	whilePreAsm += this.visit(condition)
	whilePreAsm += fmt.Sprintf("while %v @while_block_%v\n", this.Register.ReturnRegPop(), whileCounter)
	whilePreAsm += fmt.Sprintf("jump @while_end_%v\n", whileCounter)
	whilePreAsm += fmt.Sprintf("@while_block_%v\n", whileCounter)
	this.WhileCounter += 1

	blockBody := node.(parser.WhileStatement).Body
	blockAsm := this.visit(blockBody)

	var asm string
	asm += whilePreAsm
	asm += blockAsm
	asm += fmt.Sprintf("jump @while_init_%v\n", whileCounter)
	asm += fmt.Sprintf("@while_end_%v\n", whileCounter)

	return asm
}

func (this *CodeGenerator) visitForStatement(node parser.Node) string {
	if node.Type() != parser.AstTypeForStatement.Name() {
		return ""
	}

	// 先访问 init
	init := node.(parser.ForStatement).Init
	initAsm := this.visit(init)

	// 访问 test
	test := node.(parser.ForStatement).Test
	forCounter := this.ForCounter
	testAsm := fmt.Sprintf("@for_init_%v\n", forCounter)
	testAsm += this.visit(test)
	testAsm += fmt.Sprintf("for %v @for_block_%v\n", this.Register.ReturnRegPop(), forCounter)
	testAsm += fmt.Sprintf("jump @for_end_%v\n", forCounter)
	testAsm += fmt.Sprintf("@for_block_%v\n", forCounter)
	this.ForCounter += 1

	// 访问 body
	body := node.(parser.ForStatement).Body
	bodyAsm := this.visit(body)

	// 访问 update
	update := node.(parser.ForStatement).Update
	updateAsm := this.visit(update)

	// 跟 while 循环一样，到终点要调回去
	var asm string
	asm += initAsm
	asm += testAsm
	asm += bodyAsm
	asm += updateAsm
	asm += fmt.Sprintf("jump @for_init_%v\n", forCounter)
	asm += fmt.Sprintf("@for_end_%v\n", forCounter)

	return asm
}

func (this *CodeGenerator) visitArrayLiteral(node parser.Node) string {
	if node.Type() != parser.AstTypeArrayLiteral.Name() {
		return ""
	}
	var asm string
	arrayValues := node.(parser.ArrayLiteral).Values
	var vs []any
	for _, v := range arrayValues {
		vs = append(vs, this.visitRawLiteralValue(v))
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

func (this *CodeGenerator) visitFunctionExpression(node parser.Node) string {
	if node.Type() != parser.AstTypeFunctionExpression.Name() {
		return ""
	}
	var asm string

	params := node.(parser.FunctionExpression).Params
	for _, v := range params {
		if v.Type() == parser.AstTypeIdentifier.Name() {
			//             ; f1 += 2
			//            set2 a3 2
			//            add2 f1 a3 f1
			asm += fmt.Sprintf(".func_var %v\n", this.visitRawLiteralValue(v))
			this.StackOffset += 2
		}
	}

	body := node.(parser.FunctionExpression).Body
	asm += this.visit(body)

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
	var asm string
	// 做一下限制，变量名不为空
	name := node.(parser.VariableDeclaration).Name
	variableName := name.(parser.Identifier).Value
	if len(variableName) == 0 {
		utils.LogError("visitVariableDeclaration invalid left variable declaration", variableName)
		return ""
	}

	// 检测一下右边的表达式是不是函数，是的话走另外的逻辑
	right := node.(parser.VariableDeclaration).Value
	if right.Type() == parser.AstTypeFunctionExpression.Name() {
		asm = fmt.Sprintf(".function @%v\n", variableName)

		this.FunctionInfo[variableName] = &FunctionInfo{RetOffset: this.FunctionOffset, Name: variableName}
		asm += this.visitFunctionExpression(right)
		this.FunctionOffset += 2
		return asm
	}

	left := node.(parser.VariableDeclaration).Name
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		this.visit(left)
	default:
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return ""
	}

	var rightAsm string
	right = node.(parser.VariableDeclaration).Value
	rightAsm = this.visit(right)

	// 退出来寄存器
	asm = rightAsm + "push " + this.Register.ReturnRegPop() + "\n"
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

func (this *CodeGenerator) visitRawLiteralValue(node parser.Node) string {
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
	case parser.AstTypeIdentifier.Name():
		value = fmt.Sprintf("%v", node.(parser.Identifier).Value)
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

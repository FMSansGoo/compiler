package asm_vm_stack_base

import (
	"go-compiler/parser"
	"go-compiler/utils"
)

type Compiler struct {
	instructions Instructions
	constants    []Object
	scopes       []CompilationScope
	scopeIndex   int
	symbolTable  *SymbolTable
	loopIndex    int
	Loops        []Loop
}

type Loop struct {
	LoopIndex       int
	LoopBreakPos    int
	LoopContinuePos int
}

type EmittedInstruction struct {
	Opcode   OpCode
	Position int
}

type CompilationScope struct {
	instructions        Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func NewCompiler() *Compiler {
	globalScope := CompilationScope{
		instructions:        Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	for i, v := range Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		instructions: Instructions{},
		constants:    []Object{},
		// 先初始化一个空的scope
		scopes: []CompilationScope{
			globalScope,
		},
		symbolTable: symbolTable,
	}
}

func (c *Compiler) Compile(node parser.Node) {
	switch node.Type() {
	case parser.AstTypeProgram.Name():
		bodys := node.(parser.Program).Body
		for _, body := range bodys {
			c.Compile(body)
		}
		c.emit(OpCodePop)
	case parser.AstTypeBlockStatement.Name():
		body := node.(parser.BlockStatement).Body
		for _, item := range body {
			c.Compile(item)
		}
	case parser.AstTypeVariableDeclaration.Name():
		n := node.(parser.VariableDeclaration)
		name := n.Name.(parser.Identifier).Value
		utils.LogInfo("define variable", name)
		symbol := c.symbolTable.Define(name)
		c.Compile(n.Value)

		if symbol.Scope == GlobalScope {
			c.emit(OpCodeSetGlobal, symbol.Index)
		} else {
			c.emit(OpCodeSetLocal, symbol.Index)
		}
	case parser.AstTypeAssignmentExpression.Name():
		n := node.(parser.AssignmentExpression)
		// 先处理 = 的赋值
		switch n.Operator {
		case "=":
			name := n.Left.(parser.Identifier).Value
			utils.LogInfo("assign variable", name)
			symbol, ok := c.symbolTable.Resolve(name)
			if !ok {
				utils.LogError("undefined variable", name)
			}
			c.Compile(n.Right)

			if symbol.Scope == GlobalScope {
				c.emit(OpCodeSetGlobal, symbol.Index)
			} else {
				c.emit(OpCodeSetLocal, symbol.Index)
			}
		default:
			utils.LogError("unimplemented operator", n.Operator)
			//case "+=":
		}
	case parser.AstTypeIdentifier.Name():
		n := node.(parser.Identifier)
		symbol, ok := c.symbolTable.Resolve(n.Value)
		if !ok {
			utils.LogError("undefined variable", n.Value)
		}

		c.loadSymbol(symbol)
	case parser.AstTypeUnaryExpression.Name():
		n := node.(parser.UnaryExpression)
		c.Compile(n.Value)
		switch n.Operator {
		case "not":
			c.emit(OpCodeNot)
		case "-":
			c.emit(OpCodeMinus)
		default:
			utils.LogError("unknown operator", n.Operator)
		}
	case parser.AstTypeBinaryExpression.Name():
		op := node.(parser.BinaryExpression).Operator
		c.Compile(node.(parser.BinaryExpression).Left)
		c.Compile(node.(parser.BinaryExpression).Right)
		switch op {
		case "+":
			c.emit(OpCodeAdd)
		case "-":
			c.emit(OpCodeSub)
		case "*":
			c.emit(OpCodeMul)
		case "/":
			c.emit(OpCodeDiv)
		case "==":
			c.emit(OpCodeEquals)
		case "!=":
			c.emit(OpCodeNotEquals)
		case ">=":
			c.emit(OpCodeGreaterThanEquals)
		case "<=":
			c.emit(OpCodeLessThanEquals)
		case "<":
			c.emit(OpCodeLessThan)
		case ">":
			c.emit(OpCodeGreaterThan)
		case "+=":
			c.emit(OpCodeAddEquals)
		case "-=":
			c.emit(OpCodeSubEquals)
		case "*=":
			c.emit(OpCodeMulEquals)
		case "/=":
			c.emit(OpCodeDivEquals)

		default:
			utils.LogError("unknown operator", op)
		}
	case parser.AstTypeNumberLiteral.Name():
		v := node.(parser.NumberLiteral).Value
		integer := &NumberObject{Value: v}
		c.emit(OpCodeConstant, c.addConstant(integer))
	case parser.AstTypeStringLiteral.Name():
		v := node.(parser.StringLiteral).Value
		literal := &StringObject{Value: v}
		c.emit(OpCodeConstant, c.addConstant(literal))
	case parser.AstTypeBooleanLiteral.Name():
		v := node.(parser.BooleanLiteral).Value
		if v {
			c.emit(OpCodeTrue)
		} else {
			c.emit(OpCodeFalse)
		}
	case parser.AstTypeIfStatement.Name():
		n := node.(parser.IfStatement)
		condition := n.Condition
		c.Compile(condition)

		// 用 9999 当占位符
		jumpNotTruthyPos := c.emit(OpCodeJumpNotTruthy, 9999)
		c.Compile(n.Consequent)

		if c.lastInstructionIs(OpCodePop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(OpCodeJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		// 将占位符换成真实地址
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if n.Alternate == nil {
			// todo 这个 null 不知道做啥的
			//c.emit(OpCodeNull)
		} else {
			c.Compile(n.Alternate)

			if c.lastInstructionIs(OpCodePop) {
				c.removeLastPop()
			}
		}
		afterAlternative := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternative)
	case parser.AstTypeWhileStatement.Name():
		// 标记进来条件前的地址
		inLoopBeforePos := len(c.currentInstructions())
		// 要标记进入了一个循环
		c.enterLoop()

		// 先编译条件
		n := node.(parser.WhileStatement)
		condition := n.Condition
		c.Compile(condition)

		// 用 9999 当占位符,如果 condition 不是真的就跳到 while 结束
		jumpNotTruthyPos := c.emit(OpCodeJumpNotTruthy, 9999)
		c.Compile(n.Body)

		// 在这里检测有没有
		// 把 pop 去掉
		if c.lastInstructionIs(OpCodePop) {
			c.removeLastPop()
		}

		// 跳到条件编译前
		jumpPos := c.emit(OpCodeJump, 9999)

		// 将占位符换成while 结束后的地址
		afterConsequencePos := len(c.currentInstructions())

		// 处理 continue 地址
		hasContinue, continuePos := c.getContinueRetAddress()
		if hasContinue {
			c.changeOperand(continuePos, inLoopBeforePos)
		}

		// 处理 break 地址
		hasBreak, breakPos := c.getBreakRetAddress()
		if hasBreak {
			c.changeOperand(breakPos, afterConsequencePos)
		}

		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		// 将占位符换成跳到条件编译前的地址
		c.changeOperand(jumpPos, inLoopBeforePos)
		c.outLoop()
	case parser.AstTypeBreakStatement.Name():
		jumpPos := c.emit(OpCodeJump, 9999)
		c.setBreakAddress(jumpPos)
	case parser.AstTypeContinueStatement.Name():
		jumpPos := c.emit(OpCodeJump, 9999)
		c.setContinueAddress(jumpPos)
	case parser.AstTypeForStatement.Name():
		// 先编译 初始化
		n := node.(parser.ForStatement)
		init := n.Init
		if init != nil {
			c.Compile(init)
		}

		// 要标记进入了一个循环
		c.enterLoop()

		// 标记进来条件前的地址
		inLoopConditionBeforePos := len(c.currentInstructions())

		// 条件
		condition := n.Test
		c.Compile(condition)

		// 用 9999 当占位符,如果 condition 不是真的就跳到 for 结束
		jumpOutLoopNotTruthyPos := c.emit(OpCodeJumpNotTruthy, 9999)
		// jump body
		jumLoopBodyPos := c.emit(OpCodeJump, 9999)

		inLoopUpdateBeforePos := len(c.currentInstructions())

		if n.Update != nil {
			c.Compile(n.Update)
		}

		c.emit(OpCodeJump, inLoopConditionBeforePos)

		inLoopBodyBeforePos := len(c.currentInstructions())
		c.changeOperand(jumLoopBodyPos, inLoopBodyBeforePos)

		c.Compile(n.Body)

		// 在这里检测有没有 pop，把 pop 去掉
		if c.lastInstructionIs(OpCodePop) {
			c.removeLastPop()
		}

		afterConsequencePos := len(c.currentInstructions())

		// 处理 continue 地址
		hasContinue, continuePos := c.getContinueRetAddress()
		if hasContinue {
			c.changeOperand(continuePos, inLoopUpdateBeforePos)
		}

		// 处理 break 地址
		hasBreak, breakPos := c.getBreakRetAddress()
		if hasBreak {
			c.changeOperand(breakPos, afterConsequencePos)
		}

		// 跳出去循环的地址
		c.changeOperand(jumpOutLoopNotTruthyPos, afterConsequencePos)

		c.outLoop()

	case parser.AstTypeNullLiteral.Name():
		_, ok := node.(parser.NullLiteral)
		if ok {
			c.emit(OpCodeNull)
		}
	case parser.AstTypeArrayLiteral.Name():
		vs := node.(parser.ArrayLiteral).Values
		for _, v := range vs {
			c.Compile(v)
		}
		c.emit(OpCodeArray, len(vs))
	case parser.AstTypeDictLiteral.Name():
		kvs := node.(parser.DictLiteral).Values
		for _, kv := range kvs {
			k := kv.(parser.PropertyAssignment).Key
			v := kv.(parser.PropertyAssignment).Value
			c.Compile(k)
			c.Compile(v)
		}
		c.emit(OpCodeDict, len(kvs)*2)
	case parser.AstTypeFunctionExpression.Name():
		utils.LogInfo("function in?")
		c.enterScope()
		functionNode := node.(parser.FunctionExpression)

		for _, p := range functionNode.Params {
			id := p.(parser.Identifier)
			c.symbolTable.Define(id.Value)
		}
		// 这里能做处理，假设 body 没有数据，直接加上一个 null
		c.Compile(functionNode.Body)
		body := functionNode.Body
		bs := body.(parser.BlockStatement).Body
		if len(bs) == 0 {
			c.emit(OpCodeNull)
		}

		// 这里一定要 return 一个值
		if c.lastInstructionIs(OpCodePop) {
			c.replaceLastOpcode(OpCodeReturn)
		}
		if !c.lastInstructionIs(OpCodeReturn) {
			c.emit(OpCodeReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		freeSymbols := c.symbolTable.FreeSymbols
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &CompiledFunctionObject{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(functionNode.Params),
		}

		fnIndex := c.addConstant(compiledFn)
		// 闭包，第一个数函数在常量池的索引，第二个数用于指定栈中有多少自由变量需要转移到即将创建的闭包中
		c.emit(OpCodeClosure, fnIndex, len(freeSymbols))
	case parser.AstTypeCallExpression.Name():
		n := node.(parser.CallExpression)

		c.Compile(n.Object)

		for _, arg := range n.Args {
			c.Compile(arg)
		}
		c.emit(OpCodeFunctionCall, len(n.Args))
	case parser.AstTypeMemberExpression.Name():
		n := node.(parser.MemberExpression)

		if n.ElementType == "array_dict" {
			c.Compile(n.Object)
			c.Compile(n.Property)
			c.emit(OpCodeObjectCall)
		}
		// todo 支持点语法
	case parser.AstTypeReturnStatement.Name():
		v := node.(parser.ReturnStatement).Value
		if v != nil {
			c.Compile(v)
		}
		c.emit(OpCodeReturn)
	default:
		utils.LogError("unknown node type: %s", node.Type())
	}
}

func (c *Compiler) replaceLastOpcode(opcode OpCode) {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, GenerateByte(opcode))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = opcode
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) emit(op OpCode, operands ...int) int {
	ins := GenerateByte(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) setLastInstruction(op OpCode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	// 这里主要是处理闭包那套
	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) currentInstructions() Instructions {
	instructions := c.scopes[c.scopeIndex].instructions
	return instructions
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := GetOpCodeFromValue(c.currentInstructions()[opPos])
	newInstruction := GenerateByte(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) lastInstructionIs(op OpCode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

// 常量数组
func (c *Compiler) addConstant(obj Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(OpCodeGetGlobal, s.Index)
	case LocalScope:
		c.emit(OpCodeGetLocal, s.Index)
	case FreeScope:
		c.emit(OpCodeGetFree, s.Index)
	case BuiltinScope:
		c.emit(OpCodeGetBuiltin, s.Index)
	}
}

func (c *Compiler) enterLoop() {
	c.loopIndex += 1
	c.Loops = append(c.Loops, Loop{
		LoopIndex:       c.loopIndex,
		LoopBreakPos:    -1,
		LoopContinuePos: -1,
	})
}

func (c *Compiler) outLoop() {
	c.loopIndex -= 1
	c.Loops = c.Loops[:len(c.Loops)-1]
}

func (c *Compiler) setBreakAddress(addr int) {
	c.Loops[c.loopIndex-1].LoopBreakPos = addr
}

func (c *Compiler) setContinueAddress(addr int) {
	c.Loops[c.loopIndex-1].LoopContinuePos = addr
}

func (c *Compiler) getBreakRetAddress() (ok bool, retAddr int) {
	if c.Loops[c.loopIndex-1].LoopBreakPos != -1 {
		ok = true
		retAddr = c.Loops[c.loopIndex-1].LoopBreakPos
		return
	}
	return false, -1
}

func (c *Compiler) getContinueRetAddress() (ok bool, retAddr int) {
	if c.Loops[c.loopIndex-1].LoopContinuePos != -1 {
		ok = true
		retAddr = c.Loops[c.loopIndex-1].LoopContinuePos
		return
	}
	return false, -1
}

type Bytecode struct {
	Instructions Instructions
	Constants    []Object
}

func (c *Compiler) ReturnBytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

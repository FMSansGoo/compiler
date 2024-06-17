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

	return &Compiler{
		instructions: Instructions{},
		constants:    []Object{},
		// 先初始化一个空的scope
		scopes: []CompilationScope{
			globalScope,
		},
		symbolTable: NewSymbolTable(),
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
			c.emit(OpCodeGreaterThan)
		}
	case parser.AstTypeNumberLiteral.Name():
		v := node.(parser.NumberLiteral).Value
		integer := &NumberObject{Value: v}
		c.emit(OpCodeConstant, c.addConstant(integer))
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
			c.emit(OpCodeNull)
		} else {
			c.Compile(n.Alternate)

			if c.lastInstructionIs(OpCodePop) {
				c.removeLastPop()
			}
		}
		afterAlternative := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternative)
	case parser.AstTypeNullLiteral.Name():
		_, ok := node.(parser.NullLiteral)
		if ok {
			c.emit(OpCodeNull)
		}
	case parser.AstTypeStringLiteral.Name():
		v := node.(parser.StringLiteral).Value
		literal := &StringObject{Value: v}
		c.emit(OpCodeConstant, c.addConstant(literal))
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

func (c *Compiler) ReturnBytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions Instructions
	Constants    []Object
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(OpCodeGetGlobal, s.Index)
	case LocalScope:
		c.emit(OpCodeGetLocal, s.Index)
	case FreeScope:
		c.emit(OpCodeGetFree, s.Index)
	}
}

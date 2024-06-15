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
	default:
		utils.LogError("unknown node type: %s", node.Type())
	}
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
		//case LocalScope:
		//	c.emit(code.OpGetLocal, s.Index)
		//case BuiltinScope:
		//	c.emit(code.OpGetBuiltin, s.Index)
		//case FreeScope:
		//	c.emit(code.OpGetFree, s.Index)
	}
}

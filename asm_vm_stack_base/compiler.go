package asm_vm_stack_base

import (
	"fmt"
	"go-compiler/parser"
	"go-compiler/utils"
)

type Compiler struct {
	instructions Instructions
	constants    []Object
	scopes       []CompilationScope
	scopeIndex   int
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
	}
}

func (c *Compiler) Compile(node parser.Node) {
	switch node.Type() {
	case parser.AstTypeProgram.Name():
		bodys := node.(parser.Program).Body
		for _, body := range bodys {
			c.Compile(body)
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
		c.emit(OpCodePop)
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
	fmt.Printf("ins print %v\n", ins)
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
	return c.scopes[c.scopeIndex].instructions
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

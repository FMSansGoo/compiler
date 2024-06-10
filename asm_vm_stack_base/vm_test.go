package asm_vm_stack_base

import (
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"go-compiler/utils"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestVm(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"true == true", true},
		{"false == true", false},

		{"true", true},
		{"false", false},
		{"not true", false},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"-1", -1},
		{"not true", false},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		lexer := sansLexer.SansLangLexer{}
		lexer.Code = tt.input
		tokenList := lexer.TokenList()
		tokensLexer := sansLexer.TokenList{
			Tokens: tokenList,
		}
		parser := sansParser.NewSansLangParser(&tokensLexer)
		ast := parser.Parse()
		compiler := NewCompiler()
		compiler.Compile(ast)
		bytecode := compiler.ReturnBytecode()
		utils.LogInfo("bytecode: %+v", bytecode)

		vm := NewVM(bytecode)
		err := vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.GetStackTop()
		utils.LogInfo("stackElem: %+v", stackElem)
	}
}

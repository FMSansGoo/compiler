package asm_vm_stack_base

import (
	"fmt"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"go-compiler/utils"
	"testing"
)

type CompilerTest struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []CompilerTest{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeAdd),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeSub),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeMul),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 / 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeDiv),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeGreaterThan),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeEquals),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeNotEquals),
				GenerateByte(OpCodePop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestLiteral(t *testing.T) {
	tests := []CompilerTest{
		{
			input:             "1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeTrue),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeFalse),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "null",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeNull),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "not true",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeTrue),
				GenerateByte(OpCodeNot),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeMinus),
				GenerateByte(OpCodePop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestIf(t *testing.T) {
	tests := []CompilerTest{
		{
			input: `
			if (true) { 10 } 
			`,
			expectedConstants: []interface{}{10},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeTrue),
				// 10 也是地址
				GenerateByte(OpCodeJumpNotTruthy, 10),
				GenerateByte(OpCodeConstant, 0),
				// 11 也是地址
				GenerateByte(OpCodeJump, 11),
				GenerateByte(OpCodeNull),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `
			if (true) { 10 } else { 20 }
			`,
			expectedConstants: []interface{}{10, 20},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeTrue),
				// 10 也是地址
				GenerateByte(OpCodeJumpNotTruthy, 10),
				GenerateByte(OpCodeConstant, 0),
				// 13 也是地址
				GenerateByte(OpCodeJump, 13),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodePop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalVarStatements(t *testing.T) {
	tests := []CompilerTest{
		{
			input: `
				var one = 1
				one
				`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeGetGlobal, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `
				var one = 1
				var two = 2
				`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeSetGlobal, 1),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `
				var one = 1
				var two = 2
				two
				`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeSetGlobal, 1),
				GenerateByte(OpCodeGetGlobal, 1),
				GenerateByte(OpCodePop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringAndArrayAndObject(t *testing.T) {
	tests := []CompilerTest{
		{
			input:             `"monkey"`,
			expectedConstants: []interface{}{"monkey"},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []interface{}{"mon", "key"},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeAdd),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `[]`,
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeArray, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `[1,2]`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeArray, 2),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `{}`,
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeDict, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `{"0": 5 * 6}`,
			expectedConstants: []interface{}{"0", 5, 6},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeConstant, 2),
				GenerateByte(OpCodeMul),
				GenerateByte(OpCodeDict, 2),
				GenerateByte(OpCodePop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []CompilerTest) {
	for _, tt := range tests {
		fmt.Printf("--- %s ---\n", tt.input)
		lexer := sansLexer.SansLangLexer{}
		lexer.Code = tt.input
		tokenList := lexer.TokenList()
		tokensLexer := sansLexer.TokenList{
			Tokens: tokenList,
		}
		fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)
		parser := sansParser.NewSansLangParser(&tokensLexer)
		ast := parser.Parse()
		compiler := NewCompiler()
		compiler.Compile(ast)
		bytecode := compiler.ReturnBytecode()

		err := testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			utils.LogError("testInstructions failed", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)

		if err != nil {
			utils.LogError("testConstants failed", err)
		}

	}
}

func testConstants(t *testing.T, expected []interface{}, actual []Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch c := constant.(type) {
		case int:
			err := testIntegerObject(int64(c), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case string:
			err := testStringObject(c, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual Object) error {
	result, ok := actual.(*NumberObject)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", result.Value, expected)
	}
	return nil
}

func testNullObject(expected Object, actual Object) error {
	result, ok := actual.(*NullObject)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", result, expected)
	}
	return nil
}

func testStringObject(expected string, actual Object) error {
	result, ok := actual.(*StringObject)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testInstructions(expected []Instructions, actual Instructions) error {
	concatInstructions := func(ins []Instructions) Instructions {
		out := Instructions{}
		for _, i := range ins {
			out = append(out, i...)
		}
		return out
	}

	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot%q", i, concatted, actual)
		}
	}

	return nil
}

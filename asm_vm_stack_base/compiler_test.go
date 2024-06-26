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
				GenerateByte(OpCodeGreaterThanEquals),
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

func TestAssignment(t *testing.T) {
	tests := []CompilerTest{
		{
			input:             "var a = 1 a = 1",
			expectedConstants: []interface{}{1, 1},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeGetGlobal, 0),
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
				GenerateByte(OpCodeTrue), // 1
				// 10 也是地址
				GenerateByte(OpCodeJumpNotTruthy, 10), // 3
				GenerateByte(OpCodeConstant, 0),       // 3
				// 11 也是地址
				GenerateByte(OpCodeJump, 10), //3
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
			},
		},
		{
			input: `
			if (1 == 1) { 10 } else { 20 }
			`,
			expectedConstants: []interface{}{1, 1, 10, 20},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeEquals),

				// 10 也是地址
				GenerateByte(OpCodeJumpNotTruthy, 16),
				GenerateByte(OpCodeConstant, 2),
				// 13 也是地址
				GenerateByte(OpCodeJump, 19),
				GenerateByte(OpCodeConstant, 3),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestWhile(t *testing.T) {
	tests := []CompilerTest{
		{
			input: `
			var i = 0
			while (i == 0) { i = 1 }
			`,
			expectedConstants: []interface{}{0, 0, 1},
			expectedInstructions: []Instructions{
				// pre
				GenerateByte(OpCodeConstant, 0),  // 3
				GenerateByte(OpCodeSetGlobal, 0), // 6
				// conition
				GenerateByte(OpCodeGetGlobal, 0), // 9
				GenerateByte(OpCodeConstant, 1),  // 12
				GenerateByte(OpCodeEquals),       // 13
				// body
				GenerateByte(OpCodeJumpNotTruthy, 25), // 16
				GenerateByte(OpCodeConstant, 2),       // 19
				GenerateByte(OpCodeSetGlobal, 0),      // 22
				GenerateByte(OpCodeJump, 6),           // 25
			},
		},
		{
			input: `
			var i = 0
			while (i == 0) {
				i = 1
				while(true){
					if(i == 0){
						break
					} else {
						i = 0
						continue
					}
				}
				i = 2
			} i
			`,
			expectedConstants: []interface{}{0, 0, 1},
			expectedInstructions: []Instructions{
				// pre
				GenerateByte(OpCodeConstant, 0),  // 3
				GenerateByte(OpCodeSetGlobal, 0), // 6
				// conition
				GenerateByte(OpCodeGetGlobal, 0), // 9
				GenerateByte(OpCodeConstant, 1),  // 12
				GenerateByte(OpCodeEquals),       // 13
				// body
				GenerateByte(OpCodeJumpNotTruthy, 25), // 16
				GenerateByte(OpCodeConstant, 2),       // 19
				GenerateByte(OpCodeSetGlobal, 0),      // 22
				GenerateByte(OpCodeJump, 6),           // 25
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

func TestIndexExpressions(t *testing.T) {
	tests := []CompilerTest{
		{
			input:             "[1, 2, 3][1]",
			expectedConstants: []interface{}{1, 2, 3, 1},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeConstant, 2),
				GenerateByte(OpCodeArray, 3),
				GenerateByte(OpCodeConstant, 3),
				GenerateByte(OpCodeObjectCall),
				GenerateByte(OpCodePop),
			},
		},
		{
			input:             `{"key":1}["key"]`,
			expectedConstants: []interface{}{"key", 1, "key"},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeConstant, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeDict, 2),
				GenerateByte(OpCodeConstant, 2),
				GenerateByte(OpCodeObjectCall),
				GenerateByte(OpCodePop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLambdaFunction(t *testing.T) {
	tests := []CompilerTest{
		{
			input: `function() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				// 把整个函数当做常量来返回
				[]Instructions{
					GenerateByte(OpCodeConstant, 0),
					GenerateByte(OpCodeConstant, 1),
					GenerateByte(OpCodeAdd),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				//GenerateByte(OpCodeConstant, 2),
				GenerateByte(OpCodeClosure, 2, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `function() {  5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				// 把整个函数当做常量来返回
				[]Instructions{
					GenerateByte(OpCodeConstant, 0),
					GenerateByte(OpCodeConstant, 1),
					GenerateByte(OpCodeAdd),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				//GenerateByte(OpCodeConstant, 2),
				GenerateByte(OpCodeClosure, 2, 0),
				GenerateByte(OpCodePop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFunction(t *testing.T) {
	tests := []CompilerTest{
		{
			// 无参数函数
			input: `
				const fuck = function() {
					return 1
				}
				fuck()
			`,
			expectedConstants: []interface{}{
				1,
				// 把整个函数当做常量来返回
				[]Instructions{
					GenerateByte(OpCodeConstant, 0),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeClosure, 1, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeGetGlobal, 0),
				GenerateByte(OpCodeFunctionCall, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			// 无参数函数
			input: `
				const fuck = function() { 
				}
				fuck()
			`,
			expectedConstants: []interface{}{
				// 把整个函数当做常量来返回
				[]Instructions{
					GenerateByte(OpCodeNull),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeClosure, 0, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeGetGlobal, 0),
				GenerateByte(OpCodeFunctionCall, 0),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `
			const oneArg = function(a) { }
			oneArg(24)
			`,
			expectedConstants: []interface{}{
				[]Instructions{
					GenerateByte(OpCodeNull),
					GenerateByte(OpCodeReturn),
				},
				24,
			},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeClosure, 0, 0),
				GenerateByte(OpCodeSetGlobal, 0),
				GenerateByte(OpCodeGetGlobal, 0),
				GenerateByte(OpCodeConstant, 1),
				GenerateByte(OpCodeFunctionCall, 1),
				GenerateByte(OpCodePop),
			},
		},
		{
			input: `
			function() { 
				var num = 55
				return num
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]Instructions{
					GenerateByte(OpCodeConstant, 0),
					GenerateByte(OpCodeSetLocal, 0),
					GenerateByte(OpCodeGetLocal, 0),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeClosure, 1, 0),
				GenerateByte(OpCodePop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestClosure(t *testing.T) {
	tests := []CompilerTest{
		{
			// 无参数函数
			input: `
			function(a) {
				function(b) {
					return a + b
				}
			}
			`,
			expectedConstants: []interface{}{
				[]Instructions{
					GenerateByte(OpCodeGetFree, 0),
					GenerateByte(OpCodeGetLocal, 0),
					GenerateByte(OpCodeAdd),
					GenerateByte(OpCodeReturn),
				},
				[]Instructions{
					GenerateByte(OpCodeGetLocal, 0),
					GenerateByte(OpCodeClosure, 0, 1),
					GenerateByte(OpCodeReturn),
				},
			},
			expectedInstructions: []Instructions{
				GenerateByte(OpCodeClosure, 1, 0),
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

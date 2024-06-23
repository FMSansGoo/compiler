package asm_vm_stack_base

import (
	"fmt"
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
		//{"if (true) { 10 }", 10},
		//{"if (true) { 10 } else { 20 }", 10},
		//{"if (false) { 10 } else { 20 }", 20},
		//{"if (1) { 10 }", 10},
		{`var i = 0 
		while (i == 0) { 
			while(true){
				break
			} 
			i = 1
		} i`, 1,
		},
		{`var i = 0 
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
		} i`, 2,
		},
	}

	runVmTests(t, tests)
}

func TestGlobalVariableStatements(t *testing.T) {
	tests := []vmTestCase{
		{"var one = 1\none", 1},
		{"var one = 1\nvar two = 2\none + two", 3},
		{"var one = 1\nvar two = one + one\none + two", 3},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"sans"`, "sans"},
		{`"sans" + "one"`, "sansone"},
		{`"sans" + "one" + "beloved"`, "sansonebeloved"},
	}

	runVmTests(t, tests)
}

func TestArrayExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []any{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
		{`["fk"]`, []string{"fk"}},
	}
	runVmTests(t, tests)
}

func TestObjectLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[DictKeyObject]Object{},
		},
		{
			`{"1": 2}`,
			map[DictKeyObject]Object{
				DictKeyObject{Key: StringObject{Value: "1"}}: &NumberObject{Value: 2},
			},
		},
	}

	runVmTests(t, tests)
}

func TestArrayDictIndexExpression(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 2, 3]][0][0]", 1},
		{`{"k":1}["k"]`, &NumberObject{Value: 1}},
		{`
				var a = 1	
				{a:1}[a]
			`, &NumberObject{Value: 1}},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		//{
		//	input: `
		//	const test = function() { 1 }
		//	test()
		//	`,
		//	expected: 1,
		//},
		{
			input: `
			const test = function() { }
			test()
			`,
			expected: nil,
		},
		{
			input: `
		const a = function() { 1 }
		const b = function() { a() + 1 } 
		const c = function() { b() + 1 }
		c()
			`,
			expected: 3,
		},
		{
			input: `
		const a = function() { 1 }
		const b = function() { return a  } 
		b()()
			`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsReturn(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			const test = function() { 1 }
			test()
			`,
			expected: 1,
		},
		{
			input: `
			const test = function() { }
			test()
			`,
			expected: nil,
		},
		{
			input: `
		const a = function() { 1 }
		const b = function() { a() + 1 } 
		const c = function() { b() + 1 }
		c()
			`,
			expected: 3,
		},
		{
			input: `
		const a = function() { 1 }
		const b = function() { return a  } 
		b()()
			`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsLocal(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			const a = function() { 
				var num = 55
				return num
			}
			a()
			`,
			expected: 55,
		},
		{
			input: `
			const firstFoobar = function() { var foobar = 50 return foobar }
			const secondFoobar = function() { var foobar = 100 return foobar }
			firstFoobar() + secondFoobar()
			`,
			expected: 150,
		},
		{
			input: `
			var globalSeed = 50
			const minusOne = function() {
				var num = 1
				globalSeed - num
			}
			const minusTwo = function() {
				var num = 2
				globalSeed - num
			}
			minusOne() + minusTwo()
			`,
			expected: 97,
		},
		{
			input: ` 
		const returnsOneReturner = function() { 
			const returnsOne = function() { 
				return 1 
			}
			returnsOne 
		}
		returnsOneReturner()()
		`,
			expected: 1,
		},
		{
			input: ` 
		 const identity = function(a) { a }
		 identity(4)
		`,
			expected: 4,
		},
		{
			input: ` 
		 const sum = function(a, b) { return a + b }
		 sum(1,2)
		`,
			expected: 3,
		},
		{
			input: ` 
		 const sum = function(a, b) { return a + b }
		 const outer = function() { 
		 	sum(1, 2) + sum(3, 4)
		 }
		 outer()
		`,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}

func TestClosureCall(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			const newClosure = function(a) {
				return function() { a }
			}
			const closure = newClosure(99)
			closure()
			`,
			expected: 99,
		},
		{
			input: `
		const newClosure = function(a, b) {
			const one = function() { a  } 
			const two = function() { b  } 
			function() { one() + two()  } 
		} 
		const closure = newClosure(9, 90)
		closure()
		`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		 const countDown = function(x) { 
			 if (x == 0) { 
				return 0 
			 } else { 
				countDown(x - 1) 
			 } 
		 } 
		 countDown(1)
		`,
			expected: 0,
		},
		{
			input: `
			const countDown = function(x) { 
				if (x == 0) { 
					return 0  
				} else { 
					countDown(x - 1) 
				} 
			}  
			const wrapper = function() { 
				countDown(1)  
			} 
			wrapper()
		`,
			expected: 0,
		},
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

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual Object) {
	t.Helper()

	switch e := expected.(type) {
	case int:
		err := testIntegerObjectValue(int64(e), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case string:
		err := testStringObjectValue(e, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case bool:
		err := testBoolObjectValue(e, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *StringObject:
		if e.Value != actual.(*StringObject).Value {
			t.Errorf("object has wrong value. got=%q, want=%q", actual.(*StringObject).Value, e.Value)
		}
	case *NumberObject:
		if e.Value != actual.(*NumberObject).Value {
			t.Errorf("object has wrong value. got=%v, want=%v", actual.(*NumberObject).Value, e.Value)
		}
	case *BoolObject:
		if e.Value != actual.(*BoolObject).Value {
			t.Errorf("object has wrong value. got=%v, want=%v", actual.(*NumberObject).Value, e.Value)
		}
	case *NullObject:
		err := testNullObjectValue(e, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*ArrayObject)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Values) != len(e) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(e), len(array.Values))
			return
		}

		for i, expectedElm := range e {
			testExpectedObject(t, expectedElm, array.Values[i])
		}
	case []bool:
		array, ok := actual.(*ArrayObject)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Values) != len(e) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(e), len(array.Values))
			return
		}

		for i, expectedElm := range e {
			testExpectedObject(t, expectedElm, array.Values[i])
		}
	case []string:
		array, ok := actual.(*ArrayObject)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Values) != len(e) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(e), len(array.Values))
			return
		}

		for i, expectedElm := range e {
			testExpectedObject(t, expectedElm, array.Values[i])
		}
	case []any:
		array, ok := actual.(*ArrayObject)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Values) != len(e) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(e), len(array.Values))
			return
		}

		for i, expectedElm := range e {
			testExpectedObject(t, expectedElm, array.Values[i])
		}
	case map[Object]Object:
		dict, ok := actual.(*DictObject)
		if !ok {
			t.Errorf("object not Dict: %T (%+v)", actual, actual)
			return
		}
		if len(dict.Pairs) != len(e) {
			t.Errorf("hash has wrong number of Pairs. want=%d, got=%d", len(e), len(dict.Pairs))
			return
		}
		expectedKeys := make([]Object, 0)
		expectedValues := make([]Object, 0)
		for expectedKey, expectedValue := range e {
			expectedKeys = append(expectedKeys, expectedKey)
			expectedValues = append(expectedValues, expectedValue)
		}
		actualKeys := make([]Object, 0)
		actualValues := make([]Object, 0)
		for expectedKey, expectedValue := range actual.(*DictObject).Pairs {
			actualKeys = append(actualKeys, expectedKey)
			actualValues = append(actualValues, expectedValue)
		}
		for i, key := range expectedKeys {
			actualKey := actualKeys[i]
			testExpectedObject(t, key, actualKey)
		}
		for i, value := range expectedValues {
			actualValue := actualValues[i]
			testExpectedObject(t, value, actualValue)
		}
	case nil:
		result, ok := actual.(*NullObject)
		if !ok {
			t.Errorf("object is not Integer. got=%T (%+v)", result, expected)
		}
	default:
		t.Errorf("expected:%v ", e)
	}

}

func testIntegerObjectValue(expected int64, actual Object) error {
	result, ok := actual.(*NumberObject)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", result.Value, expected)
	}
	if int64(result.Value) != expected {
		return fmt.Errorf("object has wrong value. got=%+v, want=%+v",
			result.Value, expected)
	}
	return nil
}

func testNullObjectValue(expected Object, actual Object) error {
	result, ok := actual.(*NullObject)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", result, expected)
	}
	return nil
}

func testStringObjectValue(expected string, actual Object) error {
	result, ok := actual.(*StringObject)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testBoolObjectValue(expected bool, actual Object) error {
	result, ok := actual.(*BoolObject)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%+v, want=%+v", result.Value, expected)
	}

	return nil
}

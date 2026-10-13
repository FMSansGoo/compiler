package asm_vm_stack_base

import (
	"testing"

	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
)

func runSource(t *testing.T, input string, jit bool) Object {
	t.Helper()
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = input
	tokens := sansLexer.TokenList{Tokens: lexer.TokenList()}
	parser := sansParser.NewSansLangParser(&tokens)
	ast := parser.Parse()
	compiler := NewCompiler()
	compiler.Compile(ast)
	vm := NewVM(compiler.ReturnBytecode())
	vm.JitEnabled = jit
	if err := vm.Run(); err != nil {
		t.Fatalf("vm error jit=%v: %s", jit, err)
	}
	return vm.GetStackTop()
}

func TestNumericCodeMatchesInterpreter(t *testing.T) {
	programs := []struct {
		input string
		want  float64
	}{
		{
			input: `
			const f = function() { return 1 + 2 * 3 }
			f()
			`,
			want: 7,
		},
		{
			input: `
			const f = function() {
				var a = 1
				a = a + 2
				return a
			}
			f()
			`,
			want: 3,
		},
		{
			input: `
			const f = function(n) { return -n }
			f(1)
			`,
			want: -1,
		},
		{
			input: `
			const f = function(n) { return n - 1 }
			f(4)
			`,
			want: 3,
		},
	}

	for _, program := range programs {
		interpreted := runSource(t, program.input, false)
		compiled := runSource(t, program.input, true)
		if err := testIntegerObjectValue(int64(program.want), interpreted); err != nil {
			t.Fatalf("interpreter: %s\n%s", err, program.input)
		}
		if err := testIntegerObjectValue(int64(program.want), compiled); err != nil {
			t.Fatalf("numeric code: %s\n%s", err, program.input)
		}
	}
}

func TestArithmeticLowersToNumericAdd(t *testing.T) {
	input := `
	const f = function() { return 1 + 2 }
	f()
	`
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = input
	tokens := sansLexer.TokenList{Tokens: lexer.TokenList()}
	parser := sansParser.NewSansLangParser(&tokens)
	compiler := NewCompiler()
	compiler.Compile(parser.Parse())
	bytecode := compiler.ReturnBytecode()

	var fn *CompiledFunctionObject
	for _, constant := range bytecode.Constants {
		if compiled, ok := constant.(*CompiledFunctionObject); ok {
			fn = compiled
		}
	}
	if fn == nil {
		t.Fatal("compiled function missing")
	}
	code := compileNumeric(fn)
	if code == nil {
		t.Fatal("expected numeric code")
	}
	foundAdd := false
	for _, instr := range code.Instrs {
		if instr.Op == NumAdd {
			foundAdd = true
		}
	}
	if !foundAdd {
		t.Fatalf("instructions = %+v", code.Instrs)
	}
}

func TestCallIsNotNumericCode(t *testing.T) {
	input := `
	const add = function(a, b) { return a + b }
	const main = function() { return add(1, 2) }
	main()
	`
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = input
	tokens := sansLexer.TokenList{Tokens: lexer.TokenList()}
	parser := sansParser.NewSansLangParser(&tokens)
	compiler := NewCompiler()
	compiler.Compile(parser.Parse())

	var mainFn *CompiledFunctionObject
	for _, constant := range compiler.ReturnBytecode().Constants {
		compiled, ok := constant.(*CompiledFunctionObject)
		if !ok {
			continue
		}
		if compiled.NumParameters == 0 && len(compiled.Instructions) > 0 {
			mainFn = compiled
		}
	}
	if mainFn == nil {
		t.Fatal("main function missing")
	}
	if compileNumeric(mainFn) != nil {
		t.Fatal("function call should stay on the interpreter")
	}

	interpreted := runSource(t, input, false)
	compiled := runSource(t, input, true)
	if err := testIntegerObjectValue(3, interpreted); err != nil {
		t.Fatal(err)
	}
	if err := testIntegerObjectValue(3, compiled); err != nil {
		t.Fatal(err)
	}
}

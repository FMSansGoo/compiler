package repl

import (
	"bufio"
	"fmt"
	"go-compiler/asm_vm_stack_base"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"go-compiler/semantic"
	"io"
	"strings"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scope := semantic.NewScopeV2()
	symbolTable := asm_vm_stack_base.NewSymbolTable()
	for i, v := range asm_vm_stack_base.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}
	globals := make([]asm_vm_stack_base.Object, asm_vm_stack_base.GlobalSize)

	var pending strings.Builder
	for {
		if braceDepth(pending.String()) == 0 {
			fmt.Fprint(out, Prompt)
		} else {
			fmt.Fprint(out, ".. ")
		}
		if !scanner.Scan() {
			return
		}

		pending.WriteString(scanner.Text())
		pending.WriteByte('\n')
		if braceDepth(pending.String()) > 0 {
			continue
		}

		source := pending.String()
		pending.Reset()
		runSource(source, out, scope, symbolTable, globals)
	}
}

func runSource(source string, out io.Writer, scope *semantic.ScopeV2, symbolTable *asm_vm_stack_base.SymbolTable, globals []asm_vm_stack_base.Object) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Fprintf(out, "%v\n", recovered)
		}
	}()

	if strings.TrimSpace(source) == "" {
		return
	}

	l := sansLexer.NewSansLangLexer(source)
	tokensLexer := sansLexer.TokenList{
		Tokens: l.TokenList(),
	}
	p := sansParser.NewSansLangParser(&tokensLexer)
	ast := p.Parse()

	analysis := semantic.NewSemanticAnalysisV2(ast)
	analysis.CurrentScope = scope
	analysis.Visit()

	compiler := asm_vm_stack_base.NewCompilerWithState(symbolTable)
	compiler.Compile(ast)
	bytecode := compiler.ReturnBytecode()

	vm := asm_vm_stack_base.NewVMWithGlobals(bytecode, globals)
	err := vm.Run()
	if err != nil {
		fmt.Fprintf(out, "vm.Run failed: %v\n", err)
		return
	}

	stackElem := vm.GetStackTop()
	if stackElem != nil {
		io.WriteString(out, stackElem.Inspect())
	}
	io.WriteString(out, "\n")
}

func braceDepth(source string) int {
	depth := 0
	inString := byte(0)
	for i := 0; i < len(source); i++ {
		ch := source[i]
		if inString != 0 {
			if ch == '\\' && i+1 < len(source) {
				i++
				continue
			}
			if ch == inString {
				inString = 0
			}
			continue
		}
		if ch == '/' && i+1 < len(source) && source[i+1] == '/' {
			for i < len(source) && source[i] != '\n' {
				i++
			}
			continue
		}
		if ch == '"' || ch == '\'' || ch == '`' {
			inString = ch
			continue
		}
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
		}
	}
	if depth < 0 {
		return 0
	}
	return depth
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

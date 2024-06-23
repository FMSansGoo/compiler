package repl

import (
	"bufio"
	"fmt"
	"go-compiler/asm_vm_stack_base"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"go-compiler/utils"
	"io"
)

const Prompt = ">> "

// 不完整
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	symbolTable := asm_vm_stack_base.NewSymbolTable()
	for i, v := range asm_vm_stack_base.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := sansLexer.NewSansLangLexer(line)
		tokensLexer := sansLexer.TokenList{
			Tokens: l.TokenList(),
		}
		p := sansParser.NewSansLangParser(&tokensLexer)
		ast := p.Parse()

		compiler := asm_vm_stack_base.NewCompiler()
		compiler.Compile(ast)
		bytecode := compiler.ReturnBytecode()

		vm := asm_vm_stack_base.NewVM(bytecode)
		err := vm.Run()
		if err != nil {
			utils.LogErrorFormat("vm.Run failed", err)
		}

		stackElem := vm.GetStackTop()
		if stackElem != nil {
			io.WriteString(out, stackElem.Inspect())
		}
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

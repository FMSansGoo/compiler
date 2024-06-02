package asm_vm_register_base

import (
	"encoding/json"
	"fmt"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"testing"
)

func TestCodeGeneratorAsm(t *testing.T) {
	// 基本类型
	//TestCodeGeneratorAsm1_1(t)
	//TestCodeGeneratorAsm1_2(t)

	// if
	testCodeGenerator2_1()

}

func TestCodeGeneratorAsm1_1(t *testing.T) {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	1
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := sansLexer.TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
	fmt.Println("====================== parser init =======================")
	parser := sansParser.NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
	fmt.Println("====================== parser end =======================")
	fmt.Println("====================== asm_gen init =======================")

	// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)

	fmt.Println("====================== asm_gen end =======================")
	fmt.Println("====================== asm_gen_2 init =======================")

	// 生成汇编代码
	assembler := NewAssembler()
	assembler.AddAsm(codeGen.Asm)
	memory := assembler.Compile()
	fmt.Printf("memory:%+v\n", memory)

	fmt.Println("====================== asm_gen_2 end =======================")
}

func TestCodeGeneratorAsm1_2(t *testing.T) {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	1 + 1
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := sansLexer.TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
	fmt.Println("====================== parser init =======================")
	parser := sansParser.NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
	fmt.Println("====================== parser end =======================")
	fmt.Println("====================== asm_gen init =======================")

	// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)

	fmt.Println("====================== asm_gen end =======================")
	fmt.Println("====================== asm_gen_2 init =======================")

	// 生成汇编代码
	assembler := NewAssembler()
	assembler.AddAsm(codeGen.Asm)
	memory := assembler.Compile()
	fmt.Printf("memory:%+v\n", memory)

	fmt.Println("====================== asm_gen_2 end =======================")
}

func TestCodeGeneratorAsm2_1(t *testing.T) {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	1
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := sansLexer.TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
	fmt.Println("====================== parser init =======================")
	parser := sansParser.NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
	fmt.Println("====================== parser end =======================")
	fmt.Println("====================== asm_gen init =======================")

	// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)

	fmt.Println("====================== asm_gen end =======================")
	//fmt.Println("====================== asm_gen_2 init =======================")
	//
	//// 生成汇编代码
	//assembler := NewAssembler()
	//assembler.AddAsm(codeGen.Asm)
	//memory := assembler.Compile()
	//fmt.Printf("memory:%+v\n", memory)
	//
	//fmt.Println("====================== asm_gen_2 end =======================")
}

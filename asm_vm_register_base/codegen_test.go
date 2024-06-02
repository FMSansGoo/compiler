package asm_vm_register_base

import (
	"encoding/json"
	"fmt"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"testing"
)

func TestCodeGenerator(t *testing.T) {
	// 基本类型
	testCodeGenerator1_1()
	//testCodeGenerator1_2()
	//testCodeGenerator1_3()
	//testCodeGenerator1_4()
	//testCodeGenerator1_5()

	// if else 完成
	//testCodeGenerator2_1()
	//while
	//testCodeGenerator2_2()
	//for
	//testCodeGenerator2_3()

	// function
	//testCodeGenerator3_1()

}

func testCodeGenerator1_1() {
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

}

func testCodeGenerator1_2() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	1 + 1 * 1
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

}

func testCodeGenerator1_3() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	var a = 1 + 1 * 2
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

}

func testCodeGenerator1_4() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	var a = 1 
	var b = 2
	a = 3 + b
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

}

func testCodeGenerator1_5() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	var a = [1,2,3,4]
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

}

func testCodeGenerator2_1() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	if(false) {
		1
	} else if(false){
		2
	} else {
		3
	}
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

	//// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)
	//
	fmt.Println("====================== asm_gen end =======================")

}

func testCodeGenerator2_2() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
while(true){
	var a = 1
}
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

	//// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)
	//
	fmt.Println("====================== asm_gen end =======================")

}

func testCodeGenerator2_3() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	for(var a = 0; a < 1; a = a + 1) {
	}
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

	//// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)
	//
	fmt.Println("====================== asm_gen end =======================")

}

func testCodeGenerator3_1() {
	// 基本类型
	lexer := sansLexer.SansLangLexer{}
	//.call @main 1 2
	//.function @main
	// .func_var a
	// .func_var b
	// .return 0
	lexer.Code = `
	// 调用函数在定义之前
	main(1,2)
	const main = function(a,b) {
		var c = 1
		return 0
	}
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

	//// 生成汇编代码
	codeGen := NewCodeGenerator(ast)
	codeGen.Visit()
	fmt.Printf("%v\n", codeGen.Asm)
	//
	fmt.Println("====================== asm_gen end =======================")

}

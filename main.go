package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 普通运算
	test1()
	// 递归下降处理
	//test2()
	//// 优先级
	//test3()
	//// 多重表达式
	//test4()
	//// function
	//test5()
	//// if else （还没支持 else if)
	//test6()
	//// for
	//test7()
	//// while
	//test8()
	//// class
	//test9()
	//// not and or
	//test10()
	//// array、dot
	//test11()
}
func test1() {
	lexer := SansLangLexer{}
	lexer.Code = `
	var b = "1"
	const afunc = function() {
		var a = 1
		var c = 1 + b
	}
	const cfunc = function() {
		var c = 1
		var a = 1	
		var d = 1 + 1
	}
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)
	fmt.Println("====================== token end =======================")
	fmt.Println("====================== parser init =======================")
	parser := NewSansLangParser(&tokensLexer)
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
	fmt.Println("====================== NewSemanticAnalysis init =======================")
	semanticAnalysis := NewSemanticAnalysis(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")

}

func test2() {
	lexer := SansLangLexer{}
	lexer.Code = `var a = 1 + 1 * 1 - 2
	`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}

func test3() {
	lexer := SansLangLexer{}
	lexer.Code = `var a = (1+2)*2
	`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}

func test4() {
	lexer := SansLangLexer{}
	lexer.Code = `var a = 1
	var b = 1
	const c = 3`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}

func test5() {
	lexer := SansLangLexer{}
	lexer.Code = `var a = function(a, b) {
		var v = 1
	}`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}

func test6() {
	lexer := SansLangLexer{}
	lexer.Code = `if(a == 1) {
		a > 1
	} else {
		a < 2
	}`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}
func test7() {
	lexer := SansLangLexer{}
	lexer.Code = `for(var i = 0; i < 10; i+=1) {
		a = 1
	}`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}
func test8() {
	lexer := SansLangLexer{}
	lexer.Code = `while(1){
		a = 1
	}`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
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
}

func test9() {
	lexer := SansLangLexer{}
	lexer.Code = `class a super b {
	}`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)

	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
}

func test10() {
	lexer := SansLangLexer{}
	lexer.Code = `return a and b`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
}

func test11() {
	lexer := SansLangLexer{}
	lexer.Code = `a = [1,2]
	a = a.b`
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	fmt.Printf("Tokens %+v\n", tokensLexer.Tokens)
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()

	// 将节点转换为JSON字符串
	jsonData, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}

	// 打印JSON字符串
	fmt.Println(string(jsonData))
}

func test12() {

}

package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSemanticAnalysisV2(t *testing.T) {
	//基本类型
	//testSemanticAnalysis1_1()
	////多重运算
	//testSemanticAnalysis2_1()
	////for 循环 while 循环 if else continue return
	//testSemanticAnalysis3_1()
	//testSemanticAnalysis3_2()
	//testSemanticAnalysis3_3()
	////函数
	testSemanticAnalysis4_1()
	//类
	//testSemanticAnalysis5_1()
}

func testSemanticAnalysis1_1() {
	lexer := SansLangLexer{}
	lexer.Code = `
		var a = 1
		var b = "string"
		const c = true
		const d = null
		var e = [1,2,3]
		var f = {"a":1,"b":2}
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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis2_1() {
	lexer := SansLangLexer{}
	lexer.Code = `
		var a = 1 + 1 * 1 / 3 - 1
		var aa = true
		var bb = false
		var b = aa and bb
		var c = aa or b
		var d = false
		var e = not d
		var f = (1 + 1) * 2
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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis3_1() {
	lexer := SansLangLexer{}
	lexer.Code = `
		while(true){
			break
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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis3_2() {
	lexer := SansLangLexer{}
	lexer.Code = `
		for(var a = 1; a <= 10; a += 1) {

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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis3_3() {
	lexer := SansLangLexer{}
	lexer.Code = `
		var a = 1
		if(a == 1){
			continue
		} else if (a == 2) {
			return a	
		} else {
			break
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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis4_1() {
	lexer := SansLangLexer{}
	// todo 还需要做一下传参的检测
	lexer.Code = `
		const log = function(a) {
			return a
		}
		const main = function() {
			var a = 1
			a = log(a)
			return a
		}
		var b = main()
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
	semanticAnalysis := NewSemanticAnalysisV2(ast)
	semanticAnalysis.visit()
	fmt.Println("====================== NewSemanticAnalysis end =======================")
}

func testSemanticAnalysis5_1() {
	// todo 这的 class 可能有点问题，先这样吧，
	// 类函数
	lexer := SansLangLexer{}
	lexer.Code = `
		class A {
			const cls.age = 1
			const new = function() {
			}
			const cls.fuck = function() {
				var a = 1
				return a
			}
		}
		var c = A.new() // 类 的类型 ， 可以直接类似 global['A'].fuck.returnType 这样
		var d = c.fuck() // d 应该是 number类型
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

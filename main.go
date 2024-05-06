package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	test()
	// 普通运算
	//test1()
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
	// class
	testClass()
}

func test() {

}

func test1() {
	lexer := SansLangLexer{}
	lexer.Code = `
	var a = not true
	a = true
	// 这时候 a 的作用域应该是什么？
	// 如果是强类型语言要做 报错

	for(var i = 0;i < 1; i += 1) {
		a = false
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
	lexer.Code = `
	while(1){
		var a = 1
	}
	var b = 1
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

func testClass() {
	complete := `
	class A super B {
		const cls.age = 1
		
		const new = function() {
			this.gender = "boy"
		}
		const cls.fuck = function() {
		
		}
	}
	var a = A()
	a.new()
	a.fuck()
	`
	fmt.Println(complete)
	lexer := SansLangLexer{}
	lexer.Code = `
	class A {
		const cls.age = 1
		const new = function() {
		}
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

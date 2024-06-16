package parser

import (
	"fmt"
	sansLexer "go-compiler/lexer"
	"testing"
)

func TestNewBaseParser(t *testing.T) {
	//基本类型
	//testParser1()
	//多重运算
	//testParser2()
	//for 循环 while 循环 if else
	//testParser3()
	//// 函数
	//testParser4()
	//// 类
	//testParser5()
	// 额外的扩充
	testParser6()
}

func testParser1() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	var a = 1
	var b = "string"
	const c = true
	const d = null
	var e = [1,2,3]
	var f = {"a":1,"b":2}
	var d = -1
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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	// 打印JSON字符串
	fmt.Println("====================== parser end =======================")
}

func testParser2() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
		var a = 1 + 1 * 1

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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	fmt.Println("====================== parser end =======================")
}

func testParser3() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
		for(var a = 1; a <= 10; a += 1) {
			if(a == 1){
				continue
			} else if (a == 2) {
				continue	
			} else {
				break
			}
			return 
		}
		while(true){
			return
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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	fmt.Println("====================== parser end =======================")
}

func testParser4() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
		const log = function(a) {
			return a
		}
		const main = function() {
			var a = 1
			log(a)
		}
		main()
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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	fmt.Println("====================== parser end =======================")
}

func testParser5() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
		class B {
			const new = function() {
			}
		}
		class A super B {
			const cls.age = 1
			
			const new = function(name) {
				this.gender = "boy"
				this.name = name
			}
			const cls.fuck = function() {
				var a = 1
				return a
			}
		}
		var a = A.new()
		var b = a.fuck()
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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	fmt.Println("====================== parser end =======================")
}

func testParser6() {
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
		function() { return 5 + 10 }
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
	parser := NewSansLangParser(&tokensLexer)
	ast := parser.Parse()
	fmt.Printf("Ast %+v\n", ast)

	fmt.Println("====================== parser end =======================")
}

package lexer

import (
	"fmt"
	"testing"
)

func TestNewBaseLexer(t *testing.T) {
	// 基本类型
	testLexer1()
	// 多重运算
	testLexer2()
	// for 循环 while 循环 if else
	testLexer3()
	// 函数
	testLexer4()
	// 类
	testLexer5()
}

func testLexer1() {
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
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
}

func testLexer2() {
	lexer := SansLangLexer{}
	lexer.Code = `
		var a = 1 + 1 * 1 / 3 - 1
		var b = a and b
		var c = a or b
		var d = false
		var e = not d
		var f = (1 + 1) * 2
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
}

func testLexer3() {
	lexer := SansLangLexer{}
	lexer.Code = `
		for(var a = 1; a <= 10; a += 1) {
			if(a == 1){
				continue
			} else {
				break	
			}
			return 
		}
		while(true){
			break 
		}
	`
	fmt.Println("====================== token init =======================")
	tokenList := lexer.TokenList()
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
}

func testLexer4() {
	lexer := SansLangLexer{}
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
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
}

func testLexer5() {
	lexer := SansLangLexer{}
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
	tokensLexer := TokenList{
		Tokens: tokenList,
	}
	for _, token := range tokensLexer.Tokens {
		fmt.Printf("token %+v\n", token)
	}
	fmt.Println("====================== token end =======================")
}

package code_gen

import (
	"encoding/json"
	"fmt"
	sansLexer "go-compiler/lexer"
	sansParser "go-compiler/parser"
	"testing"
)

func TestCodeGenerator(t *testing.T) {
	// 基本类型
	testCodeGenerator1()
}

func testCodeGenerator1() {
	// 基本类型
	// todo 这是错的
	lexer := sansLexer.SansLangLexer{}
	lexer.Code = `
	var a = 2 * 3 + 1 * 4
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

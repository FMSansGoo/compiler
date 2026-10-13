package main

import (
	"fmt"
	"go-compiler/repl"
	"os"
	"os/user"
)

func main() {
	// 入口只启动 REPL。词法、语法、语义检查、栈式编译和执行都在 repl.runSource。
	user, err := user.Current()
	if err != nil {
		print(err)
	}

	fmt.Printf("Hello %s! This is the sans programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

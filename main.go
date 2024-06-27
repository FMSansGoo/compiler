package main

import (
	"fmt"
	"go-compiler/repl"
	"os"
	"os/user"
)

func main() {
	// 我这里暂时把 test 文件拆分了
	// 其实这里应该放整体编译的逻辑
	// 没关系先这样

	user, err := user.Current()
	if err != nil {
		print(err)
	}

	fmt.Printf("Hello %s! This is the sans programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/nicolerobin/monkey/repl"
)

func main() {
	userName, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", userName.Username)
	fmt.Println("Feel free to type in commands!")
	repl.StartVM(os.Stdin, os.Stdout)
}

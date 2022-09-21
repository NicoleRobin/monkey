package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/nicolerobin/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Println("Feel free to type in commands!")
	repl.Start(os.Stdin, os.Stdout)
}

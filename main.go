package main

import (
	"fmt"
	"github.com/nicolerobin/monkey/repl"
	"os"
	user2 "os/user"
)

func main() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Println("Feel free to type in commands!")
	repl.Start(os.Stdin, os.Stdout)
}

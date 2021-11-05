package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatalln("Invalid command call. It expected more arguments.")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}

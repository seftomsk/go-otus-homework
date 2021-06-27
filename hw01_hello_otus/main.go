package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	rawData := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(rawData))
}

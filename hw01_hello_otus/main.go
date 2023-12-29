package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	s := "Hello, OTUS!"
	fmt.Print(stringutil.Reverse(s))
}

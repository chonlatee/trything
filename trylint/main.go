package main

import (
	"fmt"
	"trylint/cal"
)

type foo struct {
	a string
}

func main() {
	_ = foo{}
	fmt.Println(cal.Add(1, 2))
}

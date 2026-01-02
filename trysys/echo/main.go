package main

import (
	"fmt"
	"os"
)

func main() {
	for i, arg := range os.Args[1:] {
		if i > 0 {
			fmt.Print(" ")
		}

		fmt.Print(arg)
	}

	fmt.Println()
}

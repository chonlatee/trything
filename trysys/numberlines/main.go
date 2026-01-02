package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: numberlines <filename>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v", filepath, err)
		os.Exit(1)
	}

	c := 0
	i := 0
	for range b {
		if b[i] == '\n' {
			c++
		}
		i++
	}

	fmt.Fprintf(os.Stdout, "%d", c)
}

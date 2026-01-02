package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: escapetext <filename>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v", filepath, err)
		os.Exit(1)
	}

	var r strings.Builder

	i := 0
	r.Grow(len(b))
	for range b {

		if b[i] == '\n' {
			r.WriteString("\\n")
		}

		if b[i] == '\r' {
			r.WriteString("\\r")
		}

		if b[i] == '\t' {
			r.WriteString("\\t")
		}

		if b[i] == '\v' {
			r.WriteString("\\v")
		}

		r.WriteString(string(b[i]))
		i++
	}

	fmt.Fprintf(os.Stdout, "%s", r.String())
}

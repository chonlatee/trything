package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	q := make(chan string, 3_000_000)

	write(q, "/home/chonlatee/go")

	close(q)

	read(q)

}

func read(q chan string) {
	for f := range q {
		fmt.Print(f)
	}
}

func write(q chan string, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("os stat err: %v\n", err)
		return
	}

	fmt.Print("d")

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			write(q, full)
		} else {
			q <- "."
		}
	}

}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

func main() {

	q := make(chan string)

	wg.Add(1)
	go read(q)
	write(q, "/home/chonlatee/Downloads")
	close(q)
	wg.Wait()

	fmt.Println()
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

func read(q chan string) {
	defer wg.Done()
	for f := range q {
		fmt.Print(f)
	}
}

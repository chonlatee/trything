package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

func main() {

	q := make(chan string, 3_000_000)

	wg.Add(1)
	write(q, "/home/chonlatee/Downloads")

	wg.Wait()

	close(q)

	read(q)

}

func write(q chan string, path string) {
	defer wg.Done()
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("os stat err: %v\n", err)
		return
	}

	fmt.Print("d")

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			wg.Add(1)
			go write(q, full)
		} else {
			q <- "."
		}
	}

}

func read(q chan string) {
	for f := range q {
		fmt.Print(f)
	}
}

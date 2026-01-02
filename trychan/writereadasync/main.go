package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var writeWG sync.WaitGroup
var readWG sync.WaitGroup

func main() {

	q := make(chan string)

	readWG.Add(1)
	go read(q)

	writeWG.Add(1)
	write(q, "/home/chonlatee/Downloads")
	writeWG.Wait()
	close(q)

	readWG.Wait()

}

func write(q chan string, path string) {
	defer writeWG.Done()
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("os stat err: %v\n", err)
		return
	}

	fmt.Print("d")
	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			writeWG.Add(1)
			go write(q, full)
		} else {
			q <- "."
		}
	}

}

func read(q chan string) {
	defer readWG.Done()
	for f := range q {
		fmt.Print(f)
	}
}

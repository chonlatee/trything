package main

import (
	"context"
	"net/http"
	"os"
)

func main() {
	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(hellWorld)}

	go server.ListenAndServe()

	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "http://localhost:8080", nil)
	resp, err := new(http.Client).Do(req)
	_ = err

	defer resp.Body.Close()

	resp.Write(os.Stdout)

}

func hellWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world\r\n"))
}

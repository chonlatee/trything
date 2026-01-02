package main

import (
	"context"
	"net/http"
	"os"
)

func main() {
	server := http.Server{Addr: ":8080", Handler: TextHandler("hello, word!")}

	go server.ListenAndServe()

	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "http://localhost:8080", nil)
	resp, err := new(http.Client).Do(req)
	_ = err

	defer resp.Body.Close()

	resp.Write(os.Stdout)

}

type TextHandler string

var _ http.Handler = TextHandler("")

func (t TextHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(t))
}

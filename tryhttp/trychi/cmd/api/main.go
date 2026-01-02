package main

import (
	"chonlatee/trychi/internal/route"
	"net/http"
)

func main() {

	router := route.Router()

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	srv.ListenAndServe()
}

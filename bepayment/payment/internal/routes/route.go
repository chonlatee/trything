package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Route() http.Handler {
	r := chi.NewRouter()

	r.Post("/generate-link", func(w http.ResponseWriter, r *http.Request) {

	})
	r.Post("/inquiry-link", func(w http.ResponseWriter, r *http.Request) {

	})

	return r
}

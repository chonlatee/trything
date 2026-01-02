package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Snake struct {
	Name         string `json:"name" binding:"required"`
	PoisionLevel *int   `json:"poisionLevel" binding:"required,gte=0"`
	Description  string `json:"description" binding:"required"`
}

var validate *validator.Validate

func shouldBindJSON[T any](r *http.Request) (T, error) {
	var req T
	v := validator.New()
	json.NewDecoder(r.Body).Decode(&req)
	v.SetTagName("binding")

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("json")

		if name == "" || name == "-" {
			return ""
		}

		return strings.Split(name, ",")[0]
	})

	err := v.Struct(req)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		log.Printf("err: %v %v", errs[0].Field(), errs[0].Tag())
		return req, errs[0]
	}

	return req, nil
}

func adminRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/insert", func(w http.ResponseWriter, r *http.Request) {

		req, err := shouldBindJSON[Snake](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("simple req: %+v", req)

		w.Write([]byte("ok"))

	})

	return r

}

func Router() http.Handler {
	v1 := chi.NewRouter()

	v1.Get("/gallery", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("simple gallery"))
	})

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})
	r.Mount("/v1", v1)
	r.Mount("/admin", adminRoute())

	return r
}

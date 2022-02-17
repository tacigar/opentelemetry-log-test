package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/baz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Tracing!"))
	})
	http.ListenAndServe(":8003", r)
}

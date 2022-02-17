package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	client := http.DefaultClient

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/foo", func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get("http://localhost:8002/bar")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(body))
	})
	http.ListenAndServe(":8001", r)
}

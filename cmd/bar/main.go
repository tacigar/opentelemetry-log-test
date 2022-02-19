package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tacigar/opentelemetry-log-test/internal/otelog"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	tp, err := otelog.InitTrace("bar", "0.0.1")
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	client := otelog.GetHTTPClient()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Method("GET", "/bar", otelog.WrapHandler(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8003/baz", http.NoBody)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(span.SpanContext().TraceID())
		w.Write([]byte(body))
	}, "bar"))

	http.ListenAndServe(":8002", r)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tacigar/opentelemetry-log-test/internal/otelog"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	tp, err := otelog.InitTrace("baz", "0.0.1")
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Method("GET", "/baz", otelog.WrapHandler(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		fmt.Println(span.SpanContext().TraceID())
		w.Write([]byte("Hello Tracing!\n"))
	}, "baz"))

	http.ListenAndServe(":8003", r)
}

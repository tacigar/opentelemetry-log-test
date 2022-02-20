package main

import (
	"context"
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

	logger, err := otelog.NewZapLogger()
	if err != nil {
		panic(err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Method("GET", "/baz", otelog.WrapHandler(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		logger.Info(otelog.LogContent{Message: "Success", Span: span})
		w.Write([]byte("Hello Tracing!\n"))
	}, "baz"))

	http.ListenAndServe(":8003", r)
}

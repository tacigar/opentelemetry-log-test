package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("serviceA"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
	)
	tracer := tracerProvider.Tracer("")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "foo")
		defer span.End()
		fmt.Println("Here")

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8002/bar", http.NoBody)
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
	}), "foo")

	r.Method("GET", "/foo", h)
	http.ListenAndServe(":8001", r)
}

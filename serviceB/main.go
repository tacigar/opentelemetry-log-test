package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func printContextInternals(ctx interface{}, inner bool, depth int) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				printContextInternals(reflectValue.Interface(), true, depth+1)
			} else {
				fmt.Printf(strings.Repeat(" ", depth)+"field name: %+v\n", reflectField.Name)
				fmt.Printf(strings.Repeat(" ", depth)+"value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("context is empty (int)\n")
	}
}

func main() {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("serviceB"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
	)
	tracer := tracerProvider.Tracer("")
	fmt.Println(tracer)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := context.Background()
		// ctx, span := tracer.Start(ctx, "bar")
		// defer span.End()

		ctx := r.Context()

		printContextInternals(ctx, true, 0)
		// fmt.Println(ctx.Value("0"))

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
		// fmt.Println(span.SpanContext().TraceID())
		w.Write([]byte(body))
	}), "bar")

	r.Method("GET", "/bar", h)
	http.ListenAndServe(":8002", r)
}

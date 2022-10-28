package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/controller"
	"app/db"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	OtelExporterOTLPEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	Port 					= os.Getenv("PORT")
	
)

func initTracer() (*sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(OtelExporterOTLPEndpoint)) 
	if err != nil {
		return nil, err
	}

	batch := sdktrace.NewBatchSpanProcessor(exporter)

	
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batch),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL,semconv.ServiceNameKey.String("APP"))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func main() {

	rand.Seed(time.Now().UnixNano())
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := mux.NewRouter()
	
	rGet := r.Methods(http.MethodGet).Subrouter()
	rGet.Use(otelmux.Middleware("server"))
	rGet.HandleFunc("/facts", controller.UselessFactsHandler)

	rHelper := r.Methods(http.MethodGet).Subrouter()
	rHelper.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.Copy(w, bytes.NewBuffer([]byte("OK")))
	})
	rHelper.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {


		err := db.PingDB()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, bytes.NewBuffer([]byte("OK")))
	})

	if Port == "" {
		Port = "8080"
	}

	server := &http.Server{
		Addr:    ":"+Port,
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 2)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down server...")

}

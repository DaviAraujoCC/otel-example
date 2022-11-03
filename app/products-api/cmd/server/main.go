package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"products-api/internal/config"
	"products-api/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func initTracer() (*sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(viper.GetString("OTEL_EXPORTER_ENDPOINT"))) 
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

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}

	config.New()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	srv := server.NewServer()
	srv.Listen()
	logrus.Println("Listening on port", viper.GetString("HTTP_PORT"))

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	srv.Shutdown()
	log.Println("Shutting down server...")

}

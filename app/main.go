package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	OtelExporterOTLPEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	ReturnError 			= os.Getenv("RETURN_ERROR")
)

func initTracer() (*sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(OtelExporterOTLPEndpoint)) 
	if err != nil {
		return nil, err
	}

	batch := sdktrace.NewBatchSpanProcessor(exporter)

	
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batch),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("ExampleService"))),
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
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()


	otelHandler := otelhttp.NewHandler(http.HandlerFunc(HelloHandler), "HelloTracer")

	http.Handle("/hello", otelHandler)
	
	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 2)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down server...")

}

func HelloHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	time.Sleep(1 * time.Second)

	message, err := GenerateRandomMessage(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ReturnError == "true" {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", message)


}

type UserInfoResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
	} `json:"results"`
}

func GenerateRandomMessage(ctx context.Context) (string, error) {

	sctx,span := otel.GetTracerProvider().Tracer("HelloTracer").Start(ctx, "ReturnMessageTracer")
	defer span.End()
	
	span.AddEvent("ReturnMessage called")

	userInfo, err := getRandomUserInfo(sctx)
	if err != nil {
		return "", err
	}

	message := "Hello " + userInfo.Results[0].Name.Title + " " + userInfo.Results[0].Name.First + " " + userInfo.Results[0].Name.Last

	time.Sleep(1 * time.Second)

	span.AddEvent("ReturnMessage finished")


	return message, nil
}

func getRandomUserInfo(ctx context.Context) (UserInfoResponse, error) {
	var userInfo UserInfoResponse

	resp, err := otelhttp.Get(ctx,"https://randomuser.me/api/")
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}
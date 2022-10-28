package controller

import (
	"app/db"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)


var ( 
	tracer = otel.Tracer("Tracer")
	returnError = os.Getenv("RETURN_ERROR")
)

func UselessFactsHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	time.Sleep(time.Duration(rand.Intn(2000-500)-500) * time.Millisecond)

	message, err := GenerateRandomFact(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if returnError == "true" {
		http.Error(w, "Return Error is true", http.StatusInternalServerError)
		return
	}

	_, span := tracer.Start(ctx, "Ping Mongo DB")

	err = db.PingDB()
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	span.End()

	time.Sleep(time.Duration(rand.Intn(2000-500)-500) * time.Millisecond)

	fmt.Fprintf(w, "%s", message)

}



func GenerateRandomFact(ctx context.Context) (string, error) {

	sctx, span := tracer.Start(ctx, "Return Random Fact")
	defer span.End()

	
	userInfo, err := getRandomFactInfo(sctx)
	if err != nil {
		return "", err
	}

	message := "Fact: " + userInfo.Text + " - Source: " + userInfo.Source


	return message, nil
}

type FactInfoResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url"`
	Language  string `json:"language"`
	Permalink string `json:"permalink"`
}

func getRandomFactInfo(ctx context.Context) (FactInfoResponse, error) {
	var userInfo FactInfoResponse

	resp, err := otelhttp.Get(ctx, "https://uselessfacts.jsph.pl/random.json?language=en")
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
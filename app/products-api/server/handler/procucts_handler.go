package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func extractQueryParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	query := r.URL.Query()
	for key, value := range query {
		params[key] = value[0]
	}
	return params
}

func ProductGetIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("id: ", id)
	w.WriteHeader(http.StatusOK)
}

func ProductGetAllHandler(w http.ResponseWriter, r *http.Request) {


	params := extractQueryParams(r)

	fmt.Println("params: ", params)

	w.WriteHeader(http.StatusOK)
}

func ProductAddHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
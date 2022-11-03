package handler

import (
	"fmt"
	"net/http"
)


func NotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed\n"
	fmt.Fprintf(w, "%s", Body)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	Body := r.URL.Path + " is not supported.\n"
	fmt.Fprintf(w, "%s", Body)
}
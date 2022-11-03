package handler

import (
	"bytes"
	"io"
	"net/http"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewBuffer([]byte("OK")))
}
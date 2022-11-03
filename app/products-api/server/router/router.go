package router

import (
	"net/http"

	"products-api/server/handler"

	"products-api/server/router/routes"

	"github.com/gorilla/mux"
)



func New() *mux.Router {
	r := mux.NewRouter()
	SetupMinimalHandlers(r)

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.NotAllowedHandler)
	r.NotFoundHandler = http.HandlerFunc(handler.DefaultHandler)

	routes.SetupRoutes(r)
	return r
}

func SetupMinimalHandlers(r *mux.Router) {
	rGet := r.Methods(http.MethodGet).Subrouter()
	rGet.HandleFunc("/healthcheck", handler.HealthcheckHandler)
	rGet.HandleFunc("/readiness", handler.ReadinessHandler)
}
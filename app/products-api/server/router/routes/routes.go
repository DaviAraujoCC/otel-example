package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type Route struct {
	Method  string
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func SetupRoutes(r *mux.Router) {

	var routes []Route

	routes = append(routes, productRoutes...)

	r.Use(otelmux.Middleware("products-api"))
	for _, route := range routes {
		r.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}
}
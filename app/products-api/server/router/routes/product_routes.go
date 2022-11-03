package routes

import (
	"net/http"
	"products-api/server/handler"
)

var productRoutes = []Route{
	{
		Path:   "/api/products/{id}",
		Method:  http.MethodGet,
		Handler: handler.ProductGetIdHandler,
	},
	{
		Path:   "/api/products",
		Method:  http.MethodGet,
		Handler: handler.ProductGetAllHandler,
	},
	{
		Path:   "/api/products",
		Method:  http.MethodPost,
		Handler: handler.ProductAddHandler,
	},
}

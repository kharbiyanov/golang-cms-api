package main

import (
	"cms-api/modules/auth"
	"cms-api/modules/products"
	"cms-api/utils"
	"github.com/graphql-go/handler"
	"net/http"
)

func main() {
	authHandler := handler.New(&handler.Config{
		Schema:   &auth.Schema,
		Pretty:   true,
		GraphiQL: false,
	})

	productsHandler := handler.New(&handler.Config{
		Schema:   &products.Schema,
		Pretty:   true,
		GraphiQL: false,
	})

	http.Handle("/auth", utils.AuthMiddleware(authHandler))
	http.Handle("/products", utils.AuthMiddleware(productsHandler))
	http.ListenAndServe(utils.Config.ServerAddr, nil)
}

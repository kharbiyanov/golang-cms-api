package main

import (
	"cms-api/modules/auth"
	"cms-api/modules/products"
	"cms-api/utils"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	serverAddr, serverAddrExist := os.LookupEnv("SERVER_ADDR")
	if !serverAddrExist {
		log.Fatal("Env SERVER_ADDR does not exist")
	}
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
	http.ListenAndServe(serverAddr, nil)
}

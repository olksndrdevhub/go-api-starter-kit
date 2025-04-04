package main

import (
	"github.com/olksndrdevhub/go-api-starter-kit/db"
	"github.com/olksndrdevhub/go-api-starter-kit/middleware"
	"log"
	"net/http"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.CloseDBConn()

	baseRouter := http.NewServeMux()

	baseRouter = loadRouters(baseRouter)

	// Create a middleware chain
	middlewareStuck := middleware.CreateStuck(
		middleware.LogingMiddleware,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: middlewareStuck(baseRouter),
	}

	log.Printf("Server is running on http://localhost:8000")
	log.Fatal(server.ListenAndServe())
}

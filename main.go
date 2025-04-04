package main

import (
	"log"
	"net/http"

	"github.com/olksndrdevhub/go-api-starter-kit/db"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.CloseDBConn()

	handler := SetupRouters()

	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	log.Printf("Server is running on http://localhost:8000")
	log.Fatal(server.ListenAndServe())
}

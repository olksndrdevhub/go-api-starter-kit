package main

import (
	"fmt"
	"github.com/oleksandrdevhub/go-api-starter-kit/db"
	"github.com/oleksandrdevhub/go-api-starter-kit/handlers"
	"github.com/oleksandrdevhub/go-api-starter-kit/middleware"
	"log"
	"net/http"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.CloseDBConn()

	// base router
	baseRouter := http.NewServeMux()

	// health check router
	statusRouter := http.NewServeMux()
	statusRouter.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok"}`)
	})

	// safe apiRouter (no auth)
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("POST /register", handlers.Register)
	apiRouter.HandleFunc("POST /login", handlers.Login)

	// unsafe API router (jwt auth)
	apiJwtRouter := http.NewServeMux()
	apiJwtRouter.HandleFunc("/me", handlers.Profile)

	// api versioning
	apiV1Router := http.NewServeMux()
	apiV1Router.Handle("/v1/auth/", http.StripPrefix("/v1/auth", apiRouter))
	apiV1Router.Handle("/v1/", middleware.JWTMiddleware(http.StripPrefix("/v1", apiJwtRouter)))

	// admin router (TODO:)
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	baseRouter.Handle("/api/", http.StripPrefix("/api", apiV1Router))
	baseRouter.Handle("/admin/", http.StripPrefix("/admin", adminRouter))
	baseRouter.Handle("/status/", http.StripPrefix("/status", statusRouter))

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

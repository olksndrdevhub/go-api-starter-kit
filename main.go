package main

import (
	"log"
	"net/http"

	"github.com/olksndrdevhub/go-api-starter-kit/db"
	"github.com/olksndrdevhub/go-api-starter-kit/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: .env file not found")
	}

	jwtSecret := utils.GetEnv("JWT_SECRET", "secret")
	utils.SetJWTSecretKey([]byte(jwtSecret))

	dbConfig := db.DBConfig{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"), // PostgreSQL default port
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "app"),
		SSLMode:  utils.GetEnv("DB_SSL_MODE", "disable"), // For PostgreSQL
	}

	// Create database instance
	err = db.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close()

	handler := SetupRouters()

	port := utils.GetEnv("PORT", "8000")
	addr := ":" + port
	log.Printf("Server is starting on %s...\n", addr)

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}

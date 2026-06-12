package main

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/routes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	log.Println("Starting application...")

	// Connect to Postgres
	config.ConnectDB()

	key := config.GenerateRandomKey()
	helpers.SetJWTKey(key)
	fmt.Printf("Generated Key: %s\n", key)

	// Build the router (pure net/http, no framework)
	handler := routes.SetupRoutes()

	// Listen address is configurable via PORT (defaults to 8080).
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

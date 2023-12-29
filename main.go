package main

import (
	"net/http"
	"os"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"CloudChat/handlers"
)

func init() {
	err := godotenv.Load()	
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
}

func main() {
	r := chi.NewRouter()
	var port = os.Getenv("API_PORT")

	// Health Check
	r.Get("/", handlers.HealthCheckHandler)

	// Registration 
	r.Post("/users", handlers.RegisterHandler)

	serverAddr := ":" + port
	log.Printf("Server is listening on %s...\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
package main

import (
	"log"
	"net/http"
	"os"

	"url-shortener/backend/handlers"
	"url-shortener/backend/persistence"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func InitServer() (*mux.Router, func(), error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables.")
	}

	// Initialize the database
	db, err := persistence.InitDB()
	if err != nil {
		return nil, nil, err
	}

	// Cleanup function to close the database connection
	cleanup := func() {
		if db != nil { // Guard against nil pointer
			db.Close()
		}
	}

	// Create a repository
	urlRepo := persistence.NewPostgresURLRepository(db)

	// Create a handlers instance
	h := handlers.NewHandlers(urlRepo)

	// Set up the router
	router := mux.NewRouter()
	handlers.RegisterRoutes(router, h)

	// Serve the frontend if enabled
	serveFrontend := os.Getenv("SERVE_FRONTEND")
	if serveFrontend == "true" {
		log.Println("Serving frontend from ./frontend/")
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))
	} else {
		log.Println("Frontend serving is disabled.")
	}

	return router, cleanup, nil
}


func main() {
	// Initialize the server
	router, cleanup, err := InitServer()
	if err != nil {
		log.Fatalf("Failed to initialize the server: %v\n", err)
	}
	defer cleanup()

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

package handlers

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, handlers *Handlers) {
	// Route for shortening a URL
	router.HandleFunc("/api/v1/shorten", handlers.ShortenURLHandler).Methods("POST")

	// Route for redirecting a short code
	router.HandleFunc("/r/{shortCode}", handlers.RedirectHandler).Methods("GET")
}

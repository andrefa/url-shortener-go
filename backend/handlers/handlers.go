package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"url-shortener/backend/persistence"

	"github.com/gorilla/mux"
)

type Handlers struct {
	UrlRepository persistence.URLRepository // Injected repository dependency
}

func NewHandlers(urlRepository persistence.URLRepository) *Handlers {
	return &Handlers{UrlRepository: urlRepository}
}

func (h *Handlers) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a unique short code
	shortCode := generateShortCode()

	// Save the URL to the database
	err := h.UrlRepository.SaveURL(shortCode, request.URL)
	if err != nil {
		log.Printf("Failed to save URL: %v\n", err)
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	host := os.Getenv("HOST")
	response := map[string]string{
		"shortenedUrl": host + "/r/" + shortCode,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Use mux.Vars to extract the short code from the URL path
	log.Println("Redirecting...")
	log.Println(r)

	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	// Retrieve the original URL from the repository
	originalURL, err := h.UrlRepository.GetOriginalURL(shortCode)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if originalURL == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

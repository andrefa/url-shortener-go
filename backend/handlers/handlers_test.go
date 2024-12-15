package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"url-shortener/backend/persistence"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test ShortenURLHandler
func TestShortenURLHandler(t *testing.T) {
	// Mock the repository
	mockRepo := new(persistence.MockURLRepository)
	mockRepo.On("SaveURL", mock.Anything, "https://example.com").Return(nil)
	os.Setenv("HOST", "http://localhost")

	// Create a handlers instance with the mock repository
	h := NewHandlers(mockRepo)

	// Create a request
	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Call the handler
	h.ShortenURLHandler(w, req)

	// Verify the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]string
	json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Contains(t, responseBody["shortenedUrl"], "http://localhost/r/")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestShortenURLHandlerInvalidBody(t *testing.T) {
	// Mock the repository (not actually used in this test)
	mockRepo := new(persistence.MockURLRepository)
	h := NewHandlers(mockRepo)
	os.Setenv("HOST", "http://localhost")

	// Create a request with an invalid body
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer([]byte("invalid-body")))
	w := httptest.NewRecorder()

	// Call the handler
	h.ShortenURLHandler(w, req)

	// Verify the response
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Test RedirectHandler
func TestRedirectHandler(t *testing.T) {
	// Mock the repository
	mockRepo := new(persistence.MockURLRepository)
	mockRepo.On("GetOriginalURL", "abc123").Return("https://example.com", nil)

	// Create a handlers instance with the mock repository
	h := NewHandlers(mockRepo)

	// Create a router and register the route
	router := mux.NewRouter()
	router.HandleFunc("/r/{shortCode}", h.RedirectHandler).Methods("GET")

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/r/abc123", nil)
	w := httptest.NewRecorder()

	// Pass the request through the router
	router.ServeHTTP(w, req)

	// Verify the response
	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com", resp.Header.Get("Location"))

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestRedirectHandlerURLNotFound(t *testing.T) {
	// Mock the repository
	mockRepo := new(persistence.MockURLRepository)
	mockRepo.On("GetOriginalURL", "not-found").Return("", nil)

	// Create a handlers instance with the mock repository
	h := NewHandlers(mockRepo)

	// Create a router and register the route
	router := mux.NewRouter()
	router.HandleFunc("/r/{shortCode}", h.RedirectHandler).Methods("GET")

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/r/not-found", nil)
	w := httptest.NewRecorder()

	// Pass the request through the router
	router.ServeHTTP(w, req)

	// Verify the response
	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestRedirectHandlerDBError(t *testing.T) {
	// Mock the repository
	mockRepo := new(persistence.MockURLRepository)
	mockRepo.On("GetOriginalURL", "abc123").Return("", assert.AnError)

	// Create a handlers instance with the mock repository
	h := NewHandlers(mockRepo)

	// Create a router and register the route
	router := mux.NewRouter()
	router.HandleFunc("/r/{shortCode}", h.RedirectHandler).Methods("GET")

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/r/abc123", nil)
	w := httptest.NewRecorder()

	// Pass the request through the router
	router.ServeHTTP(w, req)

	// Verify the response
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

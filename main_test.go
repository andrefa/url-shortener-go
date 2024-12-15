package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"url-shortener/backend/persistence"

	"github.com/stretchr/testify/assert"
)

func TestInitServer(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("SERVE_FRONTEND", "false")
	os.Setenv("PORT", "8099")

	// Mock InitDB to avoid actual database connection
	persistence.SetInitDB(func() (*sql.DB, error) {
		// Return a dummy DB for tests
		return nil, nil
	})

	// Initialize the server
	router, cleanup, err := InitServer()
	if cleanup != nil {
		defer cleanup() // Ensure the cleanup function is safe to call
	}

	// Assert no errors during server initialization
	assert.NoError(t, err, "Server should initialize without error")
	assert.NotNil(t, router, "Router should not be nil")
	assert.NotNil(t, cleanup, "Cleanup function should not be nil")

	// Test a registered route
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Unregistered route should return 404")
}

func TestFrontendServing(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("SERVE_FRONTEND", "true")
	os.Setenv("PORT", "8099")

	// Mock InitDB to avoid actual database connection
	persistence.SetInitDB(func() (*sql.DB, error) {
		// Return a dummy DB for tests
		return nil, nil
	})

	// Initialize the server
	router, cleanup, err := InitServer()
	defer cleanup()

	// Assert no errors during server initialization
	assert.NoError(t, err, "Server should initialize without error")
	assert.NotNil(t, router, "Router should not be nil")

	// Test the frontend serving
	req := httptest.NewRequest(http.MethodGet, "/somefile.html", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Frontend should return 404 for non-existent files")
}

package persistence

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveURL(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresURLRepository(db)

	// Mock the expected SQL query
	mock.ExpectExec("INSERT INTO urls").
		WithArgs("shortCode", "https://example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.SaveURL("shortCode", "https://example.com")
	assert.NoError(t, err)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOriginalURL(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresURLRepository(db)

	// Mock the expected SQL query
	mock.ExpectQuery("SELECT original_url FROM urls WHERE short_code = \\$1").
		WithArgs("shortCode").
		WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow("https://example.com"))

	// Call the method
	originalURL, err := repo.GetOriginalURL("shortCode")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOriginalURLNotFound(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresURLRepository(db)

	// Mock the expected SQL query to return no rows
	mock.ExpectQuery("SELECT original_url FROM urls WHERE short_code = \\$1").
		WithArgs("invalidCode").
		WillReturnRows(sqlmock.NewRows([]string{}))

	// Call the method
	originalURL, err := repo.GetOriginalURL("invalidCode")
	assert.EqualError(t, err, "sql: no rows in result set")
	assert.Empty(t, originalURL)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

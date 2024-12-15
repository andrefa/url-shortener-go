package persistence

import "database/sql"

// URLRepository defines the interface for the repository
type URLRepository interface {
	SaveURL(shortCode, originalURL string) error
	GetOriginalURL(shortCode string) (string, error)
}

// PostgresURLRepository is the real implementation of URLRepository
type PostgresURLRepository struct {
	DB *sql.DB // Use *sql.DB here
}

// NewPostgresURLRepository creates a new PostgresURLRepository
func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{DB: db}
}

func (r *PostgresURLRepository) SaveURL(shortCode, originalURL string) error {
	query := "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, shortCode, originalURL) // Use Exec for INSERT queries
	return err
}

func (r *PostgresURLRepository) GetOriginalURL(shortCode string) (string, error) {
	query := "SELECT original_url FROM urls WHERE short_code = $1"
	row := r.DB.QueryRow(query, shortCode) // Use QueryRow for SELECT queries
	var originalURL string
	err := row.Scan(&originalURL) // Scan the result into the originalURL variable
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

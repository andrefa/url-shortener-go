package persistence

import "github.com/stretchr/testify/mock"

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) SaveURL(shortCode, originalURL string) error {
	args := m.Called(shortCode, originalURL)
	return args.Error(0)
}

func (m *MockURLRepository) GetOriginalURL(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}

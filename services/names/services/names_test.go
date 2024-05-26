package services_test

import (
	"errors"
	"micro-names/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockNamesRepository is a mock implementation of NamesRepository for testing purposes
type MockNamesRepository struct{}

func (m *MockNamesRepository) CreateName(name string) (string, error) {
	if name == "error" {
		return "", errors.New("repository error")
	}
	return "12345", nil
}

func TestCreateName(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		expectedID    string
		expectedError error
	}{
		{
			name:          "successful creation",
			inputName:     "testName",
			expectedID:    "12345",
			expectedError: nil,
		},
		{
			name:          "error from repository",
			inputName:     "error",
			expectedID:    "",
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockNamesRepository{}
			service := services.NewNamesService(mockRepo)

			id, err := service.CreateName(tt.inputName)

			assert.Equal(t, tt.expectedID, id, "unexpected ID returned")
			assert.Equal(t, tt.expectedError, err, "unexpected error returned")
		})
	}
}

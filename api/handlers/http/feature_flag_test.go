package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"ff/internal/db/model"
	featureflag "ff/internal/feature_flag"
	"ff/internal/feature_flag/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of FeatureFlagRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddFeatureFlag(flag model.FeatureFlag) error {
	args := m.Called(flag)
	return args.Error(0)
}

var validFeatureFlagBody = entity.FeatureFlag{
	Name:        "TEST_FLAG_NAME",
	Description: "Test Description",
	IsActive:    true,
}

// ... existing imports and mock setup ...

func TestCreateFeatureFlagHandler_Success(t *testing.T) {
	tests := []struct {
		name                string
		input               entity.FeatureFlag
		personId            string
		mockRepositoryError error
		expectedStatusCode  int
		expectedResponse    string
	}{
		{
			name:               "Success",
			input:              validFeatureFlagBody,
			personId:           "123",
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"message":"Feature Flag Created"}`,
		},
	}

	runTests(t, tests, nil)
}

func TestCreateFeatureFlagHandler_InvalidInput(t *testing.T) {
	tests := []struct {
		name                string
		input               entity.FeatureFlag
		personId            string
		mockRepositoryError error
		expectedStatusCode  int
		expectedResponse    string
	}{
		{
			name:               "Name is required",
			input:              entity.FeatureFlag{},
			personId:           "123",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"name is required"}`,
		},
		{
			name:               "Name must be uppercase and contain only letters, numbers, underscores",
			input:              entity.FeatureFlag{Name: "test_flag_name"},
			personId:           "123",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"name must be uppercase and contain only letters, numbers, underscores"}`,
		},
		{
			name:               "Description is required",
			input:              entity.FeatureFlag{Name: validFeatureFlagBody.Name},
			personId:           "123",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"description is required"}`,
		},
		{
			name: "ExpirationDate format is invalid",
			input: entity.FeatureFlag{
				Name:           validFeatureFlagBody.Name,
				Description:    validFeatureFlagBody.Description,
				ExpirationDate: "invalid-date",
			},
			personId:           "123",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"expirationDate must be in YYYY-MM-DD format"}`,
		},
	}

	runTests(t, tests, nil)
}

func TestCreateFeatureFlagHandler_InternalValidationErrors(t *testing.T) {
	tests := []struct {
		name                string
		input               entity.FeatureFlag
		personId            string
		mockRepositoryError error
		expectedStatusCode  int
		expectedResponse    string
	}{
		{
			name:                "Missing PersonId",
			input:               validFeatureFlagBody,
			personId:            "",
			mockRepositoryError: nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponse:    `{"error":"missing Personid header"}`,
		},
		{
			name:                "Invalid PersonId",
			input:               validFeatureFlagBody,
			personId:            "invalid",
			mockRepositoryError: nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponse:    `{"error":"invalid Personid format"}`,
		},
		{
			name:                "Repository Error",
			input:               validFeatureFlagBody,
			personId:            "123",
			mockRepositoryError: errors.New("repository error"),
			expectedStatusCode:  http.StatusInternalServerError,
			expectedResponse:    `{"error":"repository error"}`,
		},
	}

	runTests(t, tests, nil)
}

// Helper function to run the tests
func runTests(t *testing.T, tests []struct {
	name                string
	input               entity.FeatureFlag
	personId            string
	mockRepositoryError error
	expectedStatusCode  int
	expectedResponse    string
}, setupMock func(*MockRepository)) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Personid", tt.personId)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create mock repository
			mockRepository := new(MockRepository)
			// Only set up the mock expectation if we're testing a case that should reach the repository
			if tt.expectedStatusCode == http.StatusCreated || tt.mockRepositoryError != nil {
				mockRepository.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(tt.mockRepositoryError)
			}

			if setupMock != nil {
				setupMock(mockRepository)
			}

			// Create logger
			mockLogger := zerolog.New(os.Stdout)

			// Create handler
			handler := &EchoHandler{
				FeatureFlagService: featureflag.FeatureFlagService{
					Repository: mockRepository,
					Logger:     &mockLogger,
				},
			}

			// Prepare input
			inputJSON, _ := json.Marshal(tt.input)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

			// Perform request
			err := handler.createFeatureFlagHandler(c)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tt.expectedResponse, rec.Body.String())

			// Verify mock expectations
			if tt.expectedStatusCode == http.StatusCreated || tt.mockRepositoryError != nil {
				mockRepository.AssertExpectations(t)
			}
		})
	}
}

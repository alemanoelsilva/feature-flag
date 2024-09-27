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

func (m *MockRepository) GetFeatureFlag(filters *model.FeatureFlagFilters) ([]model.FeatureFlag, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.FeatureFlag), args.Error(1)
}

// Create Feature Flag Tests Cases
func TestCreateFeatureFlagHandler(t *testing.T) {
	validFeatureFlagBody := entity.FeatureFlag{
		Name:        "TEST_FLAG_NAME",
		Description: "Test Description",
		IsActive:    true,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{}, nil)
		mockRepository.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(nil)

		mockLogger := zerolog.New(os.Stdout)

		handler := &EchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(validFeatureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Perform request
		err := handler.createFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"message":"Feature Flag Created"}`, rec.Body.String())

		mockRepository.AssertCalled(t, "GetFeatureFlag", mock.Anything)
		mockRepository.AssertCalled(t, "AddFeatureFlag", mock.Anything)
	})

	t.Run("Missing PersonId", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		// Intentionally not setting PersonId header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockLogger := zerolog.New(os.Stdout)

		handler := &EchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(validFeatureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Perform request
		err := handler.createFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"missing Personid header"}`, rec.Body.String())

		mockRepository.AssertNotCalled(t, "GetFeatureFlag")
		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
	})

	t.Run("Invalid PersonId", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "invalid")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockLogger := zerolog.New(os.Stdout)

		handler := &EchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(validFeatureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Perform request
		err := handler.createFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"invalid Personid format"}`, rec.Body.String())

		mockRepository.AssertNotCalled(t, "GetFeatureFlag")
		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
	})

	t.Run("Add Repository Error", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{}, nil)
		mockRepository.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(errors.New("add repository error"))

		mockLogger := zerolog.New(os.Stdout)

		handler := &EchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(validFeatureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Perform request
		err := handler.createFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"add repository error"}`, rec.Body.String())

		mockRepository.AssertCalled(t, "GetFeatureFlag", mock.Anything)
		mockRepository.AssertCalled(t, "AddFeatureFlag", mock.Anything)
	})

	t.Run("Get Repository Error", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{}, errors.New("get repository error"))

		mockLogger := zerolog.New(os.Stdout)

		handler := &EchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(validFeatureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Perform request
		err := handler.createFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"get repository error"}`, rec.Body.String())

		mockRepository.AssertCalled(t, "GetFeatureFlag", mock.Anything)
		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
	})
}

// Get Feature Flag Tests Cases

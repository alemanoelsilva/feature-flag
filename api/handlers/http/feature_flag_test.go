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

func (m *MockRepository) GetFeatureFlag(filters model.FeatureFlagFilters, pagination model.Pagination) ([]model.FeatureFlag, int64, error) {
	args := m.Called(filters, pagination)
	return args.Get(0).([]model.FeatureFlag), int64(args.Get(1).(int)), args.Error(2)
}

func (m *MockRepository) UpdateFeatureFlagById(id uint, featureFlag model.UpdateFeatureFlag) error {
	args := m.Called(id, featureFlag)
	return args.Error(0)
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

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.Name == "TEST_FLAG_NAME"
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})
		featureFlagMock := mock.AnythingOfType("model.FeatureFlag")

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)
		mockRepository.On("AddFeatureFlag", featureFlagMock).Return(nil)

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

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
		mockRepository.AssertCalled(t, "AddFeatureFlag", featureFlagMock)
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

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.Name == "TEST_FLAG_NAME"
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})
		featureFlagMock := mock.AnythingOfType("model.FeatureFlag")

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)
		mockRepository.On("AddFeatureFlag", featureFlagMock).Return(errors.New("add repository error"))

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

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
		mockRepository.AssertCalled(t, "AddFeatureFlag", featureFlagMock)
	})

	t.Run("Get Repository Error", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.Name == "TEST_FLAG_NAME"
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, errors.New("get repository error"))

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

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
	})
}

// Get Feature Flag Tests Cases
func TestGetFeatureFlagHandler(t *testing.T) {
	validFeatureFlagBody := entity.FeatureFlag{
		Name:        "TEST_FLAG_NAME",
		Description: "Test Description",
		IsActive:    true,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.Name == "TEST_FLAG_NAME"
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})
		featureFlagMock := mock.AnythingOfType("model.FeatureFlag")

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)
		mockRepository.On("AddFeatureFlag", featureFlagMock).Return(nil)

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

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
		mockRepository.AssertCalled(t, "AddFeatureFlag", featureFlagMock)
	})
}

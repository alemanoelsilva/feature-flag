package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"ff/internal/db/model"
	featureflag "ff/internal/feature_flag"
	featureFlagEntity "ff/internal/feature_flag/entity"
	personEntity "ff/internal/person/entity"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

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
	validFeatureFlagBody := featureFlagEntity.FeatureFlag{
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

		handler := &FeatureFlagEchoHandler{
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

		handler := &FeatureFlagEchoHandler{
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

		handler := &FeatureFlagEchoHandler{
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

	t.Run("PersonId zero (not logged)", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/feature-flags/v1/feature-flags", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "0")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := new(MockRepository)
		mockLogger := zerolog.New(os.Stdout)

		handler := &FeatureFlagEchoHandler{
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
		assert.JSONEq(t, `{"error":"you are not logged in"}`, rec.Body.String())

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

		handler := &FeatureFlagEchoHandler{
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

		handler := &FeatureFlagEchoHandler{
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
	queryString := map[string]interface{}{
		"name":  "TEST_FLAG_NAME",
		"limit": 10,
		"page":  1,
	}

	// Convert queryString to URL query parameters
	queryParams := url.Values{}
	for key, value := range queryString {
		queryParams.Add(key, fmt.Sprintf("%v", value))
	}

	timeMock := time.Now()

	featureFlagRepositoryMock := []model.FeatureFlag{{
		ID:             4,
		Name:           "TEST_FLAG_NAME",
		Description:    "This is an example feature flag",
		IsActive:       false,
		ExpirationDate: "2024-05-09",
		CreatedAt:      timeMock,
		UpdatedAt:      timeMock,
		Person: &model.Person{
			ID:    1,
			Name:  "Person Name",
			Email: "person.email@email.com",
		},
	}}

	featureFlagResponse := []featureFlagEntity.FeatureFlagResponse{{
		ID:             4,
		Name:           "TEST_FLAG_NAME",
		Description:    "This is an example feature flag",
		IsActive:       false,
		ExpirationDate: "2024-05-09",
		CreatedAt:      timeMock.Format("2006-01-02 15:04:05"),
		UpdatedAt:      timeMock.Format("2006-01-02 15:04:05"),
		Person: personEntity.PersonResponse{
			ID:    1,
			Name:  "Person Name",
			Email: "person.email@email.com",
		},
	}}

	// Create a []interface{} slice
	featureFlagsInterface := make([]interface{}, len(featureFlagResponse))

	// Loop through featureFlags and assign each element to the []interface{} slice
	for i, flag := range featureFlagResponse {
		featureFlagsInterface[i] = flag
	}

	response := PaginationResponse{
		Items: featureFlagsInterface,
		Total: 1,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/feature-flags/v1/feature-flags?"+queryParams.Encode(), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.Name == "TEST_FLAG_NAME"
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 10
		})

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return(featureFlagRepositoryMock, 1, nil)

		mockLogger := zerolog.New(os.Stdout)

		handler := &FeatureFlagEchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		// Perform request
		err := handler.getFeatureFlagHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Convert rec.Body to a map (for dynamic JSON)
		var responseBody PaginationResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &responseBody); err != nil {
			fmt.Println("Error decoding JSON:", err)
		}

		assert.Equal(t, response.Total, responseBody.Total)
		// TODO: check how to cast interface{} to model.FeatureFlagResponse
		// assert.Equal(t, featureFlagResponse[0].Name, responseBody.Items[0].Name)

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
	})
}

// Update Feature Flag By ID Tests Cases
func TestUpdateFeatureFlagByIdHandler(t *testing.T) {
	featureFlagBody := featureFlagEntity.UpdateFeatureFlag{
		Description:    "This is an example feature flag",
		IsActive:       false,
		ExpirationDate: "2024-05-09",
	}

	t.Run("Success", func(t *testing.T) {
		featureFlagId := 1
		// Setup
		e := echo.New()
		url := fmt.Sprintf("/api/feature-flags/v1/feature-flags/%d", featureFlagId)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.ID == uint(featureFlagId)
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})

		featureFlagToUpdate := model.UpdateFeatureFlag{
			Description:    "This is an example feature flag",
			IsActive:       false,
			ExpirationDate: "2024-05-09",
		}

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 1, nil)
		mockRepository.On("UpdateFeatureFlagById", uint(featureFlagId), featureFlagToUpdate).Return(nil)

		mockLogger := zerolog.New(os.Stdout)

		handler := &FeatureFlagEchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(featureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Set the path parameters, for example setting the `:id` parameter
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", featureFlagId))

		// Perform request
		err := handler.updateFeatureFlagByIdHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message":"Feature Flag Updated"}`, rec.Body.String())

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)
		mockRepository.AssertCalled(t, "UpdateFeatureFlagById", uint(featureFlagId), featureFlagToUpdate)

		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
		mockRepository.AssertNotCalled(t, "GetFeatureFlag")
	})

	t.Run("Feature flag not found", func(t *testing.T) {
		featureFlagId := 1
		// Setup
		e := echo.New()
		url := fmt.Sprintf("/api/feature-flags/v1/feature-flags/%d", featureFlagId)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Personid", "123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		filtersMock := mock.MatchedBy(func(filters model.FeatureFlagFilters) bool {
			return filters.ID == uint(featureFlagId)
		})
		paginationMock := mock.MatchedBy(func(pagination model.Pagination) bool {
			return pagination.Page == 1 && pagination.Limit == 1
		})

		mockRepository := new(MockRepository)
		mockRepository.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)

		mockLogger := zerolog.New(os.Stdout)

		handler := &FeatureFlagEchoHandler{
			FeatureFlagService: featureflag.FeatureFlagService{
				Repository: mockRepository,
				Logger:     &mockLogger,
			},
		}

		inputJSON, _ := json.Marshal(featureFlagBody)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(inputJSON))

		// Set the path parameters, for example setting the `:id` parameter
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", featureFlagId))

		// Perform request
		err := handler.updateFeatureFlagByIdHandler(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error":"feature flag not found"}`, rec.Body.String())

		mockRepository.AssertCalled(t, "GetFeatureFlag", filtersMock, paginationMock)

		mockRepository.AssertNotCalled(t, "AddFeatureFlag")
		mockRepository.AssertNotCalled(t, "GetFeatureFlag")
	})
}

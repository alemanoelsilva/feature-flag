package featureflag

import (
	"os"
	"testing"
	"time"

	"ff/internal/db/model"
	"ff/internal/feature_flag/entity"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of SqlRepository
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
func TestCreateFeatureFlag(t *testing.T) {
	t.Run("Successfully create feature flag", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		request := entity.FeatureFlag{
			Name:        "TEST_FLAG_V1",
			Description: "Test Description",
			IsActive:    true,
		}

		filtersMock := mock.AnythingOfType("model.FeatureFlagFilters")
		paginationMock := mock.AnythingOfType("model.Pagination")
		featureFlagMock := mock.AnythingOfType("model.FeatureFlag")

		mockRepo.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)
		mockRepo.On("AddFeatureFlag", featureFlagMock).Return(nil)

		err := service.CreateFeatureFlag(request, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create feature flag with expiration date", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		expirationDate := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
		request := entity.FeatureFlag{
			Name:           "TEST_FLAG_V1",
			Description:    "Test Description",
			IsActive:       true,
			ExpirationDate: expirationDate,
		}

		filtersMock := mock.AnythingOfType("model.FeatureFlagFilters")
		paginationMock := mock.AnythingOfType("model.Pagination")
		featureFlagMock := mock.AnythingOfType("model.FeatureFlag")

		mockRepo.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 0, nil)
		mockRepo.On("AddFeatureFlag", featureFlagMock).Return(nil)
		err := service.CreateFeatureFlag(request, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Duplicate feature flag", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		expirationDate := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
		request := entity.FeatureFlag{
			Name:           "TEST_DUP_FLAG_V1",
			Description:    "Test Description",
			IsActive:       true,
			ExpirationDate: expirationDate,
		}

		filtersMock := mock.AnythingOfType("model.FeatureFlagFilters")
		paginationMock := mock.AnythingOfType("model.Pagination")

		mockRepo.On("GetFeatureFlag", filtersMock, paginationMock).Return([]model.FeatureFlag{}, 1, nil)

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "feature flag already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid feature flag name", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		request := entity.FeatureFlag{
			Name:        "",
			Description: "Test Description",
			IsActive:    true,
		}

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "name is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid feature flag name", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		request := entity.FeatureFlag{
			Name:        "test flag v1",
			Description: "Test Description",
			IsActive:    true,
		}

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "name must be uppercase and contain only letters, numbers, underscores", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid feature flag description", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		request := entity.FeatureFlag{
			Name:        "TEST_FLAG_V1",
			Description: "",
			IsActive:    true,
		}

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "description is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid expiration date format", func(t *testing.T) {
		mockRepo := new(MockRepository)
		logger := zerolog.New(os.Stdout)
		service := LoadService(mockRepo, &logger)

		request := entity.FeatureFlag{
			Name:           "TEST_FLAG_V1",
			Description:    "Test Description",
			IsActive:       true,
			ExpirationDate: "invalid-date",
		}

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "expirationDate must be in YYYY-MM-DD format", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

// Get Feature Flag Tests Cases

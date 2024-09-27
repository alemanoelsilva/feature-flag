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

func (m *MockRepository) GetFeatureFlag(filters *model.FeatureFlagFilters) ([]model.FeatureFlag, error) {
	args := m.Called(filters)
	return args.Get(0).([]model.FeatureFlag), args.Error(1)
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

		mockRepo.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{}, nil)
		mockRepo.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(nil)

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

		mockRepo.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{}, nil)
		mockRepo.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(nil)

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

		mockRepo.On("GetFeatureFlag", mock.AnythingOfType("*model.FeatureFlagFilters")).Return([]model.FeatureFlag{{
			Name:        request.Name,
			Description: request.Description,
			IsActive:    request.IsActive,
		}}, nil)

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

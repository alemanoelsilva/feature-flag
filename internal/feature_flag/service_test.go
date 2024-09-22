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

// Add other necessary mock methods here...

func TestCreateFeatureFlag(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := zerolog.New(os.Stdout)
	service := LoadService(mockRepo, &logger)

	t.Run("Successfully create feature flag", func(t *testing.T) {
		request := entity.FeatureFlag{
			Name:        "Test Flag",
			Description: "Test Description",
			IsActive:    true,
		}

		mockRepo.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(nil)

		err := service.CreateFeatureFlag(request, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create feature flag with expiration date", func(t *testing.T) {
		expirationDate := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
		request := entity.FeatureFlag{
			Name:           "Test Flag with Expiration",
			Description:    "Test Description",
			IsActive:       true,
			ExpirationDate: expirationDate,
		}

		mockRepo.On("AddFeatureFlag", mock.AnythingOfType("model.FeatureFlag")).Return(nil)

		err := service.CreateFeatureFlag(request, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid expiration date format", func(t *testing.T) {
		request := entity.FeatureFlag{
			Name:           "Test Flag",
			Description:    "Test Description",
			IsActive:       true,
			ExpirationDate: "invalid-date",
		}

		err := service.CreateFeatureFlag(request, 1)

		assert.Error(t, err)
		assert.Equal(t, "invalid expiration date format", err.Error())
	})
}

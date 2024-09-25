package repository

import (
	model "ff/internal/db/model"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSqlRepository is a test suite for SqlRepository
type TestSqlRepository struct {
	suite.Suite
	db     *gorm.DB
	logger *zerolog.Logger
	repo   *SqlRepository
}

func TestFeatureFlagRepository(t *testing.T) {
	suite.Run(t, new(TestSqlRepository))
}

func (s *TestSqlRepository) SetupTest() {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)

	// Run migrations
	err = db.AutoMigrate(&model.FeatureFlag{}, &model.Person{})
	s.Require().NoError(err)

	// Create a test logger
	logger := zerolog.New(os.Stdout)

	// Initialize the repository
	s.db = db
	s.logger = &logger
	s.repo = &SqlRepository{DB: db, Logger: &logger}
}

func (s *TestSqlRepository) TearDownTest() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

func (s *TestSqlRepository) TestAddFeatureFlag() {
	// Test case 1: Successfully add a feature flag
	s.Run("Successfully add feature flag", func() {
		featureFlag := model.FeatureFlag{
			Name:        "Test Flag",
			Description: "Test Description",
			IsActive:    true,
			Person: model.Person{
				ID: 1,
			},
			PersonId: 1,
		}

		err := s.repo.AddFeatureFlag(featureFlag)
		s.Require().NoError(err)

		// Verify the feature flag was added
		var savedFlag model.FeatureFlag
		result := s.db.First(&savedFlag, "name = ?", featureFlag.Name)
		s.Require().NoError(result.Error)
		s.Equal(featureFlag.Name, savedFlag.Name)
		s.Equal(featureFlag.Description, savedFlag.Description)
		s.Equal(featureFlag.IsActive, savedFlag.IsActive)
		// s.Equal(featureFlag.Person.ID, savedFlag.Person.ID)
	})

	// Test case 2: Attempt to add a duplicate feature flag
	s.Run("Fail to add duplicate feature flag", func() {
		featureFlag := model.FeatureFlag{
			Name:        "Duplicate Flag",
			Description: "Duplicate Description",
			IsActive:    true,
			Person: model.Person{
				ID: 1,
			},
		}

		// Add the feature flag for the first time
		err := s.repo.AddFeatureFlag(featureFlag)
		s.Require().NoError(err)

		// Attempt to add the same feature flag again
		err = s.repo.AddFeatureFlag(featureFlag)
		s.Require().Error(err)
		s.Equal("error when creating feature flag", err.Error())
	})
}

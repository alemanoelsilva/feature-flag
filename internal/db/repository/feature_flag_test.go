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

// Create Feature Flag Tests Cases
func (s *TestSqlRepository) TestAddFeatureFlag() {
	// Test case 1: Successfully add a feature flag
	s.Run("Successfully add feature flag", func() {
		featureFlag := model.FeatureFlag{
			Name:        "Test Flag",
			Description: "Test Description",
			IsActive:    true,
			PersonID:    1,
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
			PersonID:    1,
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

// Get Feature Flag Tests Cases
func (s *TestSqlRepository) TestGetFeatureFlag() {
	// Test case 1: Get feature flags with a specific name
	s.Run("Get feature flags without filters", func() {
		personOnDB := model.Person{
			Name:  "Test Person",
			Email: "test@example.com",
		}
		s.db.Create(&personOnDB)
		db := s.db.First(&personOnDB, "name = ?", personOnDB.Name)
		s.Require().NoError(db.Error)

		featureFlagsOnDB := []model.FeatureFlag{
			{
				Name:        "Test Flag 1",
				Description: "Test Description 1",
				IsActive:    true,
				PersonID:    personOnDB.ID,
			},
			{
				Name:        "Test Flag 2",
				Description: "Test Description 2",
				IsActive:    false,
				PersonID:    personOnDB.ID,
			},
			{
				Name:        "Test Flag 3",
				Description: "Test Description 3",
				IsActive:    true,
				PersonID:    personOnDB.ID,
			}}

		s.db.CreateInBatches(featureFlagsOnDB, len(featureFlagsOnDB))

		var filters model.FeatureFlagFilters
		featureFlags, err := s.repo.GetFeatureFlag(&filters)
		s.Require().NoError(err)
		s.Require().Equal(len(featureFlagsOnDB), len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[0].Name, featureFlags[0].Name)
		s.Require().Equal(featureFlagsOnDB[1].Name, featureFlags[1].Name)
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[2].Name)
	})

	// Test case 2: Get feature flags with aan specific name
	s.Run("Get feature flags without filters", func() {
		personOnDB := model.Person{
			Name:  "Test Person",
			Email: "test@example.com",
		}
		s.db.Create(&personOnDB)
		db := s.db.First(&personOnDB, "name = ?", personOnDB.Name)
		s.Require().NoError(db.Error)

		featureFlagsOnDB := []model.FeatureFlag{
			{
				Name:        "Test Flag 1",
				Description: "Test Description 1",
				IsActive:    true,
				PersonID:    personOnDB.ID,
			},
			{
				Name:        "Test Flag 2",
				Description: "Test Description 2",
				IsActive:    false,
				PersonID:    personOnDB.ID,
			},
			{
				Name:        "Test Flag 3",
				Description: "Test Description 3",
				IsActive:    true,
				PersonID:    personOnDB.ID,
			}}

		s.db.CreateInBatches(featureFlagsOnDB, len(featureFlagsOnDB))

		var filters model.FeatureFlagFilters
		filters.Name = featureFlagsOnDB[2].Name
		featureFlags, err := s.repo.GetFeatureFlag(&filters)
		s.Require().NoError(err)
		s.Require().Equal(1, len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[0].Name)
	})
}

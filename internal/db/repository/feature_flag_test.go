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

var personOnDB []model.Person

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

	// insert person as pre requisite
	personOnDB = []model.Person{
		{
			Name:  "Test Person 1",
			Email: "test@example.com",
		},
		{
			Name:  "Test Person 2",
			Email: "test@example.com",
		},
	}
	s.db.Debug().CreateInBatches(personOnDB, len(personOnDB))
	db = s.db.First(&personOnDB[0], "name = ?", personOnDB[0].Name)
	s.Require().NoError(db.Error)
	db = s.db.First(&personOnDB[1], "name = ?", personOnDB[1].Name)
	s.Require().NoError(db.Error)
}

func (s *TestSqlRepository) TearDownTest() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

// Create Feature Flag Tests Cases
func (s *TestSqlRepository) TestAddFeatureFlag() {
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
	featureFlagsOnDB := []model.FeatureFlag{
		{
			Name:        "Test Flag 1",
			Description: "Test Description 1",
			IsActive:    true,
			PersonID:    personOnDB[0].ID,
		},
		{
			Name:        "Test Flag 2",
			Description: "Test Description 2",
			IsActive:    false,
			PersonID:    personOnDB[0].ID,
		},
		{
			Name:        "Test Flag 3",
			Description: "Test Description 3",
			IsActive:    true,
			PersonID:    personOnDB[1].ID,
		}}

	s.db.CreateInBatches(featureFlagsOnDB, len(featureFlagsOnDB))
	var featureFlagTest []model.FeatureFlag
	db := s.db.Debug().Find(&featureFlagTest)
	s.Require().NoError(db.Error)
	s.Require().Equal(len(featureFlagsOnDB), len(featureFlagTest))

	s.Run("Get feature flags without filters", func() {
		var filters model.FeatureFlagFilters
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 10

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(len(featureFlagsOnDB), len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[0].Name, featureFlags[0].Name)
		s.Require().Equal(featureFlagsOnDB[1].Name, featureFlags[1].Name)
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[2].Name)
		s.Require().Equal(totalCount, int64(3))
	})

	s.Run("Get feature flags without filters (page 1)", func() {
		var filters model.FeatureFlagFilters
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 1

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(len(featureFlags), 1)
		s.Require().Equal(featureFlagsOnDB[0].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(3))
	})

	s.Run("Get feature flags without filters (page 2)", func() {
		var filters model.FeatureFlagFilters
		var pagination model.Pagination
		pagination.Page = 2
		pagination.Limit = 1

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(len(featureFlags), 1)
		s.Require().Equal(featureFlagsOnDB[1].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(3))
	})

	s.Run("Get feature flags without filters (page 3)", func() {
		var filters model.FeatureFlagFilters
		var pagination model.Pagination
		pagination.Page = 3
		pagination.Limit = 1

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(len(featureFlags), 1)
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(3))
	})

	s.Run("Get feature flags with name filter", func() {
		var filters model.FeatureFlagFilters
		filters.Name = featureFlagsOnDB[2].Name
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 10

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(1, len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(1))
	})

	s.Run("Get feature flags with isActive filter", func() {
		var filters model.FeatureFlagFilters
		active := true
		filters.IsActive = &active
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 10

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(2, len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[0].Name, featureFlags[0].Name)
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[1].Name)
		s.Require().Equal(totalCount, int64(2))
	})

	s.Run("Get feature flags with id filter", func() {
		var filters model.FeatureFlagFilters
		filters.ID = 1
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 10

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(1, len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[0].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(1))
	})

	s.Run("Get feature flags with personId filter", func() {
		var filters model.FeatureFlagFilters
		filters.PersonID = personOnDB[1].ID
		var pagination model.Pagination
		pagination.Page = 1
		pagination.Limit = 10

		featureFlags, totalCount, err := s.repo.GetFeatureFlag(filters, pagination)
		s.Require().NoError(err)
		s.Require().Equal(1, len(featureFlags))
		s.Require().Equal(featureFlagsOnDB[2].Name, featureFlags[0].Name)
		s.Require().Equal(totalCount, int64(1))
	})
}

// Update Feature Flag By Id Tests Cases
func (s *TestSqlRepository) TestUpdateFeatureFlagById() {
	s.Run("Successfully update feature flag by id", func() {
		featureFlag := model.FeatureFlag{
			Name:        "Test Flag",
			Description: "Test Description",
			IsActive:    true,
			PersonID:    1,
		}

		err := s.repo.AddFeatureFlag(featureFlag)
		s.Require().NoError(err)

		var featureFlagOnDB model.FeatureFlag
		result := s.db.First(&featureFlagOnDB, "name = ?", featureFlag.Name)
		s.Require().NoError(result.Error)

		updatedFeatureFlag := model.UpdateFeatureFlag{
			Description:    featureFlag.Description,
			IsActive:       false,
			ExpirationDate: featureFlag.ExpirationDate,
		}
		err = s.repo.UpdateFeatureFlagById(featureFlagOnDB.ID, updatedFeatureFlag)
		s.Require().NoError(err)

		var savedFlag model.FeatureFlag
		result = s.db.First(&savedFlag, "id = ?", featureFlagOnDB.ID)
		s.Require().NoError(result.Error)
		s.Equal(featureFlag.Name, savedFlag.Name)
		s.Equal(featureFlag.Description, savedFlag.Description)
		// Verify the feature flag was updated
		s.Equal(updatedFeatureFlag.IsActive, savedFlag.IsActive)
		// s.Equal(featureFlag.Person.ID, savedFlag.Person.ID)
	})

	s.Run("Should no update a non existing feature flag", func() {
		updatedFeatureFlag := model.UpdateFeatureFlag{
			Description:    "Description",
			IsActive:       false,
			ExpirationDate: "2024-10-10",
		}
		err := s.repo.UpdateFeatureFlagById(99, updatedFeatureFlag)
		s.Require().Error(err)
		s.Equal("no feature flag updated", err.Error())
	})
}

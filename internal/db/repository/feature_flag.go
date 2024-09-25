package repository

import (
	"errors"
	model "ff/internal/db/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type FeatureFlagRepository interface {
	AddFeatureFlag(featureFlag model.FeatureFlag) error
	GetFeatureFlag(filters *model.FeatureFlagFilters) ([]model.FeatureFlag, error)
}

type SqlRepository struct {
	DB     *gorm.DB
	Logger *zerolog.Logger
}

func (s *SqlRepository) AddFeatureFlag(featureFlag model.FeatureFlag) error {
	if result := s.DB.Create(&featureFlag); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("error when creating feature flag")
	}

	return nil
}

func (s *SqlRepository) GetFeatureFlag(filters *model.FeatureFlagFilters) ([]model.FeatureFlag, error) {
	query := s.DB.Debug().InnerJoins("Person")
	if filters.Name != "" {
		query.Where("feature_flags.name = ?", filters.Name)
	}

	var featureFlags []model.FeatureFlag
	if result := query.Find(&featureFlags); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return nil, errors.New("error when getting feature flags")
	}

	return featureFlags, nil
}

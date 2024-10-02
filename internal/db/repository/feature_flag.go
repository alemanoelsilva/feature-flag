package repository

import (
	"errors"
	model "ff/internal/db/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type FeatureFlagRepository interface {
	AddFeatureFlag(featureFlag model.FeatureFlag) error
	GetFeatureFlag(filters model.FeatureFlagFilters, pagination model.Pagination) ([]model.FeatureFlag, int64, error)
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

func (s *SqlRepository) GetFeatureFlag(filters model.FeatureFlagFilters, pagination model.Pagination) ([]model.FeatureFlag, int64, error) {
	query := s.DB.Debug().Model(&model.FeatureFlag{}).InnerJoins("Person")

	// apply filters
	if filters.Name != "" {
		query.Where("feature_flags.name = ?", filters.Name)
	}

	// get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query.Offset(offset).Limit(pagination.Limit)

	// get feature flags
	var featureFlags []model.FeatureFlag
	if result := query.Find(&featureFlags); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return nil, 0, errors.New("error when getting feature flags")
	}

	return featureFlags, totalCount, nil
}

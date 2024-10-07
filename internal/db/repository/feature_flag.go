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
	UpdateFeatureFlagById(id uint, featureFlag model.UpdateFeatureFlag) error
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

	// this is an optional filter, it can be true/false or not be sent
	if filters.IsActive != nil {
		query.Where("feature_flags.is_active = ?", *filters.IsActive)
	}

	if filters.ID != 0 {
		query.Where("feature_flags.id = ?", filters.ID)
	}

	if filters.PersonID != 0 {
		query.Where("feature_flags.person_id = ?", filters.PersonID)
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

func (s *SqlRepository) UpdateFeatureFlagById(id uint, featureFlag model.UpdateFeatureFlag) error {
	updateData := map[string]interface{}{
		"description":     featureFlag.Description,
		"is_active":       featureFlag.IsActive, // Explicitly include even if false
		"expiration_date": featureFlag.ExpirationDate,
	}

	result := s.DB.Debug().
		Model(&model.UpdateFeatureFlag{}). // Use an empty struct for the model
		Where("id = ?", id).
		Updates(updateData)

	// result := s.DB.Debug().Model(&featureFlag).Where("id = ?", id).Updates(&featureFlag)
	if result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("error when updating feature flag")
	}
	// TODO: RowsAffected can be NotFound or Model with no updated values
	if result.RowsAffected == 0 {
		return errors.New("no feature flag updated")
	}

	return nil
}

package repository

import (
	"errors"
	model "ff/internal/db/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type FeatureFlagRepository interface {
	AddFeatureFlag(featureFlag model.FeatureFlag) error
}

type SqlRepository struct {
	DB     *gorm.DB
	Logger *zerolog.Logger
}

func (s *SqlRepository) AddFeatureFlag(featureFlag model.FeatureFlag) error {
	if result := s.DB.Create(&featureFlag); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("internal Server Error")
	}

	return nil
}

package repository

import (
	"errors"
	model "ff/internal/db/model"
	entity "ff/internal/feature-flag/entity"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SqlRepository struct {
	DB     *gorm.DB
	Logger *zerolog.Logger
}

func (s *SqlRepository) AddFeatureFlag(ff entity.FeatureFlag, personId int) error {
	expirationDate, err := time.Parse(time.DateOnly, ff.ExpirationDate)
	if err != nil {
		return errors.New("invalid expiration date format")
	}

	featureFlag := model.FeatureFlag{
		ID:             ff.ID,
		Name:           ff.Name,
		Description:    ff.Description,
		IsActive:       ff.IsActive,
		ExpirationDate: expirationDate,
		Person: model.Person{
			ID: personId,
		},
	}

	if result := s.DB.Create(&featureFlag); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("internal Server Error")
	}

	return nil
}

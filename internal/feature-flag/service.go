package featureflag

import (
	"errors"
	"ff/internal/db/model"
	repo "ff/internal/db/repository"
	entity "ff/internal/feature-flag/entity"
	"time"

	"github.com/rs/zerolog"
)

type FeatureFlagService struct {
	Repository repo.FeatureFlagRepository
	Logger     *zerolog.Logger
}

func LoadService(r repo.FeatureFlagRepository, l *zerolog.Logger) *FeatureFlagService {
	return &FeatureFlagService{
		Logger:     l,
		Repository: r,
	}
}

func (ff *FeatureFlagService) CreateFeatureFlag(request entity.FeatureFlag, personId int) error {
	ff.Logger.Info().Msg("Creating a new Feature Flag")

	var expirationDate *time.Time
	if request.ExpirationDate != "" {
		expDate, err := time.Parse(time.DateOnly, request.ExpirationDate)
		if err != nil {
			return errors.New("invalid expiration date format")
		}
		expirationDate = &expDate
	}

	return ff.Repository.AddFeatureFlag(model.FeatureFlag{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		IsActive:    request.IsActive,
		// to make this field optional, we need to pass a pointer to the expiration date
		ExpirationDate: expirationDate,
		Person: model.Person{
			ID: personId,
		},
	})
}

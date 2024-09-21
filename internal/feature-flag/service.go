package featureflag

import (
	repo "ff/internal/db/repository"
	entity "ff/internal/feature-flag/entity"

	"github.com/rs/zerolog"
)

type FeatureFlagService struct {
	Repository repo.SqlRepository
	Logger     *zerolog.Logger
}

func LoadService(r repo.SqlRepository, l *zerolog.Logger) *FeatureFlagService {
	return &FeatureFlagService{
		Logger:     l,
		Repository: r,
	}
}

func (ser *FeatureFlagService) CreateFeatureFlag(input entity.FeatureFlag, personId int) error {
	ser.Logger.Info().Msg("Creating a new Feature Flag")

	return ser.Repository.AddFeatureFlag(entity.FeatureFlag{
		ID:             input.ID,
		Name:           input.Name,
		Description:    input.Description,
		IsActive:       input.IsActive,
		ExpirationDate: input.ExpirationDate,
	}, personId)
}

package featureflag

import (
	"errors"

	"ff/internal/db/model"
	repo "ff/internal/db/repository"
	entity "ff/internal/feature_flag/entity"

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

func (ff *FeatureFlagService) CreateFeatureFlag(request entity.FeatureFlag, personId uint) error {
	ff.Logger.Info().Msg("Creating a new Feature Flag")

	if err := request.Validate(); err != nil {
		return errors.New(err.Error())
	}

	var filters model.FeatureFlagFilters
	filters.Name = request.Name

	featureFlags, err := ff.Repository.GetFeatureFlag(&filters)
	if err != nil {
		return err
	}

	// TODO: return 409
	if len(featureFlags) > 0 {
		return errors.New("feature flag already exists")
	}

	return ff.Repository.AddFeatureFlag(model.FeatureFlag{
		ID:             request.ID,
		Name:           request.Name,
		Description:    request.Description,
		IsActive:       request.IsActive,
		ExpirationDate: request.ExpirationDate,
		PersonID:       personId,
	})
}

func (ff *FeatureFlagService) GetFeatureFlag() ([]entity.FeatureFlagResponse, error) {
	ff.Logger.Info().Msg("Getting Feature Flag")

	var filters model.FeatureFlagFilters
	featureFlags, err := ff.Repository.GetFeatureFlag(&filters)
	if err != nil {
		return nil, err
	}

	var featureFlagResponses []entity.FeatureFlagResponse
	for _, ffDB := range featureFlags {
		featureFlagResponses = append(featureFlagResponses, entity.FeatureFlagResponse{
			ID:             ffDB.ID,
			Name:           ffDB.Name,
			Description:    ffDB.Description,
			IsActive:       ffDB.IsActive,
			ExpirationDate: ffDB.ExpirationDate,
			CreatedAt:      ffDB.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      ffDB.UpdatedAt.Format("2006-01-02 15:04:05"),
			Person: entity.PersonResponse{
				ID:    ffDB.Person.ID,
				Name:  ffDB.Person.Name,
				Email: ffDB.Person.Email,
			},
		})
	}

	return featureFlagResponses, nil
}

package featureflag

import (
	"errors"

	"ff/internal/db/model"
	repo "ff/internal/db/repository"
	featureFlagEntity "ff/internal/feature_flag/entity"
	personEntity "ff/internal/person/entity"

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

func (ffs *FeatureFlagService) CreateFeatureFlag(request featureFlagEntity.FeatureFlag, personId uint) error {
	ffs.Logger.Info().Msg("Creating a new Feature Flag")

	if err := request.Validate(); err != nil {
		return errors.New(err.Error())
	}

	_, totalCount, err := ffs.Repository.GetFeatureFlag(model.FeatureFlagFilters{
		Name: request.Name,
	}, model.Pagination{
		Limit: 1,
		Page:  1,
	})
	if err != nil {
		return err
	}

	if totalCount > 0 {
		return errors.New("feature flag already exists")
	}

	return ffs.Repository.AddFeatureFlag(model.FeatureFlag{
		ID:             request.ID,
		Name:           request.Name,
		Description:    request.Description,
		IsActive:       request.IsActive,
		IsGlobal:       request.IsGlobal,
		ExpirationDate: request.ExpirationDate,
		PersonID:       personId,
	})
}

func (ffs *FeatureFlagService) GetFeatureFlag(page int, limit int, name string, isActive *bool, isGlobal *bool, id uint, personId uint) ([]featureFlagEntity.FeatureFlagResponse, int64, error) {
	ffs.Logger.Info().Msg("Getting Feature Flag")

	var pagination model.Pagination
	pagination.Page = page
	pagination.Limit = limit

	var filters model.FeatureFlagFilters
	filters.Name = name
	filters.IsActive = isActive
	filters.IsGlobal = isGlobal
	filters.ID = id
	filters.PersonID = personId

	featureFlags, totalCount, err := ffs.Repository.GetFeatureFlag(filters, pagination)
	if err != nil {
		return nil, 0, err
	}

	var featureFlagResponses []featureFlagEntity.FeatureFlagResponse
	for _, ffDB := range featureFlags {
		featureFlagResponses = append(featureFlagResponses, featureFlagEntity.FeatureFlagResponse{
			ID:             ffDB.ID,
			Name:           ffDB.Name,
			Description:    ffDB.Description,
			IsActive:       ffDB.IsActive,
			IsGlobal:       ffDB.IsGlobal,
			ExpirationDate: ffDB.ExpirationDate,
			CreatedAt:      ffDB.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      ffDB.UpdatedAt.Format("2006-01-02 15:04:05"),
			Person: personEntity.PersonResponse{
				ID:    ffDB.Person.ID,
				Name:  ffDB.Person.Name,
				Email: ffDB.Person.Email,
			},
		})
	}

	return featureFlagResponses, totalCount, nil
}

func (ffs *FeatureFlagService) UpdateFeatureFlagById(id uint, request featureFlagEntity.UpdateFeatureFlag) error {
	ffs.Logger.Info().Msg("Updating a Feature Flag")

	if err := request.Validate(); err != nil {
		return errors.New(err.Error())
	}

	_, countTotal, err := ffs.Repository.GetFeatureFlag(model.FeatureFlagFilters{
		ID: id,
	}, model.Pagination{
		Limit: 1,
		Page:  1,
	})
	if err != nil {
		return errors.New(err.Error())
	}

	if countTotal == 0 {
		return errors.New("feature flag not found")
	}

	return ffs.Repository.UpdateFeatureFlagById(id, model.UpdateFeatureFlag{
		Description:    request.Description,
		IsActive:       request.IsActive,
		IsGlobal:       request.IsGlobal,
		ExpirationDate: request.ExpirationDate,
	})
}

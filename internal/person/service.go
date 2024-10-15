package featureflag

import (
	"ff/internal/db/model"
	repo "ff/internal/db/repository"
	featureFlagEntity "ff/internal/feature_flag/entity"
	personEntity "ff/internal/person/entity"

	"github.com/rs/zerolog"
)

type PeopleService struct {
	Repository repo.PersonRepository
	Logger     *zerolog.Logger
}

func LoadService(r repo.PersonRepository, l *zerolog.Logger) *PeopleService {
	return &PeopleService{
		Logger:     l,
		Repository: r,
	}
}

func (ps *PeopleService) GetPeople(page, limit int, name string) ([]personEntity.PersonResponse, int64, error) {
	ps.Logger.Info().Msg("Getting people")

	var pagination model.Pagination
	pagination.Page = page
	pagination.Limit = limit

	people, totalCount, err := ps.Repository.GetPeople(pagination, name)
	if err != nil {
		return nil, 0, err
	}

	var personResponses []personEntity.PersonResponse
	for _, pDB := range people {
		personResponses = append(personResponses, personEntity.PersonResponse{
			ID:    pDB.ID,
			Name:  pDB.Name,
			Email: pDB.Email,
		})
	}

	return personResponses, totalCount, nil
}

func (ps *PeopleService) GetPeopleAssignmentByFeatureFlag(page, limit int, id uint, name string, isAssigned *bool) ([]personEntity.PersonWithAssignmentResponse, int64, error) {
	ps.Logger.Info().Msg("Getting people w/ assignment")

	var pagination model.Pagination
	pagination.Page = page
	pagination.Limit = limit

	people, totalCount, err := ps.Repository.GetPeopleAssignmentByFeatureFlag(pagination, id, name, isAssigned)
	if err != nil {
		return nil, 0, err
	}

	var personResponses []personEntity.PersonWithAssignmentResponse
	for _, pDB := range people {
		personResponses = append(personResponses, personEntity.PersonWithAssignmentResponse{
			ID:         pDB.ID,
			Name:       pDB.Name,
			Email:      pDB.Email,
			IsAssigned: pDB.IsGlobal,
		})
	}

	return personResponses, totalCount, nil
}

func (ps *PeopleService) GetAssignedFeatureFlagsByPersonId(id uint) ([]featureFlagEntity.AssignedFeatureFlagResponse, error) {
	ps.Logger.Info().Msg("Getting assigned feature flags by person id")

	featureFlags, err := ps.Repository.GetAssignedFeatureFlagsByPersonId(id)
	if err != nil {
		return nil, err
	}

	var featureFlagResponses []featureFlagEntity.AssignedFeatureFlagResponse
	for _, ffDB := range featureFlags {
		if ffDB.IsAssigned || ffDB.IsGlobal {
			featureFlagResponses = append(featureFlagResponses, featureFlagEntity.AssignedFeatureFlagResponse{
				ID:         ffDB.ID,
				Name:       ffDB.Name,
				IsAssigned: ffDB.IsGlobal || ffDB.IsAssigned,
			})
		}

	}

	return featureFlagResponses, nil
}

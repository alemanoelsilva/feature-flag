package person

import (
	"ff/internal/db/model"
	p_entity "ff/internal/person/entity"
	"strconv"

	"github.com/rs/zerolog"
)

type PersonRepository interface {
	GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, filters p_entity.PersonFilters) ([]model.PersonWithAssignment, int64, error)
	GetAssignedFeatureFlagsByPersonId(id uint) ([]model.AssignedFeatureFlag, error)
}

type PeopleService struct {
	Repository PersonRepository
	Logger     *zerolog.Logger
}

func LoadService(r PersonRepository, l *zerolog.Logger) *PeopleService {
	return &PeopleService{
		Logger:     l,
		Repository: r,
	}
}

func (ps *PeopleService) GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, filters p_entity.PersonFilters) ([]p_entity.PersonWithAssignmentResponse, int64, error) {
	ps.Logger.Info().Msg("Getting people w/ assignment")

	people, totalCount, err := ps.Repository.GetPeopleAssignmentByFeatureFlag(pagination, filters)
	if err != nil {
		return nil, 0, err
	}

	var personResponses []p_entity.PersonWithAssignmentResponse
	for _, pDB := range people {
		personResponses = append(personResponses, p_entity.PersonWithAssignmentResponse{
			ID:         strconv.Itoa(int(pDB.ID)),
			Name:       pDB.Name,
			Email:      pDB.Email,
			IsAssigned: pDB.IsGlobal,
		})
	}

	return personResponses, totalCount, nil
}

func (ps *PeopleService) GetAssignedFeatureFlagsByPersonId(id uint) ([]p_entity.AssignedFeatureFlagResponse, error) {
	ps.Logger.Info().Msg("Getting assigned feature flags by person id")

	featureFlags, err := ps.Repository.GetAssignedFeatureFlagsByPersonId(id)
	if err != nil {
		return nil, err
	}

	var featureFlagResponses []p_entity.AssignedFeatureFlagResponse
	for _, ffDB := range featureFlags {
		if ffDB.IsAssigned || ffDB.IsGlobal {
			featureFlagResponses = append(featureFlagResponses, p_entity.AssignedFeatureFlagResponse{
				ID:         ffDB.ID,
				Name:       ffDB.Name,
				IsActive:   ffDB.IsActive,
				IsAssigned: ffDB.IsGlobal || ffDB.IsAssigned,
			})
		}

	}

	return featureFlagResponses, nil
}

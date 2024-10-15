package repository

import (
	"errors"
	model "ff/internal/db/model"
)

type PersonRepository interface {
	GetPeople(pagination model.Pagination, name string) ([]model.Person, int64, error)
	GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, id uint, name string, isAssigned *bool) ([]model.PersonWithAssignment, int64, error)
	GetAssignedFeatureFlagsByPersonId(id uint) ([]model.AssignedFeatureFlag, error)
}

func (s *SqlRepository) GetPeople(pagination model.Pagination, name string) ([]model.Person, int64, error) {
	query := s.DB.Debug().Model(&model.Person{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query.Offset(offset).Limit(pagination.Limit)

	var people []model.Person
	if result := query.Find(&people); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return nil, 0, errors.New("error when getting people")
	}

	return people, totalCount, nil
}

func (s *SqlRepository) GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, id uint, name string, isAssigned *bool) ([]model.PersonWithAssignment, int64, error) {
	var featureFlag model.FeatureFlag
	err := s.DB.Debug().Model(&model.FeatureFlag{}).Where("id = ?", id).Find(&featureFlag).Error
	if err != nil {
		s.Logger.Error().Err(err)
		return nil, 0, errors.New("error when getting feature flag")
	}

	query := s.DB.Debug().
		Table("person p").
		Select("p.id, p.name, p.email, IF(ffa.id IS NULL, false, true) AS is_assigned").
		Joins("LEFT JOIN feature_flag_assignments ffa ON ffa.person_id = p.id AND ffa.feature_flag_id = ?", id).
		Joins("LEFT JOIN feature_flags ff ON ff.id = ffa.feature_flag_id").
		Order("p.id")

	if name != "" {
		query = query.Where("p.name LIKE ?", "%"+name+"%")
	}

	if isAssigned != nil && *isAssigned {
		query.Where("(ff.is_global = false AND ffa.id IS NOT NULL)")
	}

	// get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query.Offset(offset).Limit(pagination.Limit)

	var people []model.PersonWithAssignment
	err = query.Scan(&people).Error
	if err != nil {
		s.Logger.Error().Err(err)
		return nil, 0, errors.New("error when getting people")
	}

	var response []model.PersonWithAssignment
	for _, person := range people {
		response = append(response, model.PersonWithAssignment{
			ID:         person.ID,
			Name:       person.Name,
			Email:      person.Email,
			IsAssigned: person.IsAssigned,
			IsGlobal:   featureFlag.IsGlobal || person.IsAssigned,
		})
	}

	return response, totalCount, nil
}

func (s *SqlRepository) GetAssignedFeatureFlagsByPersonId(id uint) ([]model.AssignedFeatureFlag, error) {
	var featureFlags []model.AssignedFeatureFlag

	// TODO: handle feature flag global
	err := s.DB.Debug().Model(&model.AssignedFeatureFlag{}).Table("feature_flags ff").
		Select("ff.id, ff.name, ff.is_active, ff.is_global, if(ffa.id is null, false, true) is_assigned ").
		Joins("LEFT JOIN feature_flag_assignments ffa ON ffa.feature_flag_id = ff.id AND ffa.person_id = ?", id).
		Order("ff.id").
		Scan(&featureFlags).Error

	if err != nil {
		s.Logger.Error().Err(err)
		return nil, errors.New("error when getting feature flag by people")
	}

	return featureFlags, nil
}

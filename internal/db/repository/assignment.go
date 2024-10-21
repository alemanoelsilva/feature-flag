package repository

import (
	"errors"
	model "ff/internal/db/model"
)

func (s *SqlRepository) ApplyAssignment(assignment model.Assignment) error {
	if result := s.DB.Debug().Create(&assignment); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("error when assigning feature flag")
	}

	return nil
}

func (s *SqlRepository) GetAssignmentsByPersonAndFeatureFlagId(personId, featureFlagId uint) (model.Assignment, error) {
	query := s.DB.Debug().Model(&model.Assignment{}).Where("person_id = ?", personId).Where("feature_flag_id = ?", featureFlagId)

	// get feature flags
	var assignment model.Assignment
	if result := query.Find(&assignment); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return model.Assignment{}, errors.New("error when getting assignment")
	}

	return assignment, nil
}

func (s *SqlRepository) DeleteAssignment(assignment model.Assignment) error {
	if result := s.DB.Debug().Where("person_id = ? AND feature_flag_id = ?", assignment.PersonID, assignment.FeatureFlagID).Delete(&model.Assignment{}); result.Error != nil {
		s.Logger.Error().Err(result.Error)
		return errors.New("error when deleting assigning")
	}

	return nil
}

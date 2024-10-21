package assignment

import (
	"errors"
	"fmt"

	assignmentEntity "ff/internal/assignment/entity"
	"ff/internal/db/model"

	"github.com/rs/zerolog"
)

type AssignmentRepository interface {
	ApplyAssignment(assignment model.Assignment) error
	GetAssignmentsByPersonAndFeatureFlagId(personId, featureFlagId uint) (model.Assignment, error)
	DeleteAssignment(assignment model.Assignment) error
}

type AssignmentService struct {
	Repository AssignmentRepository
	Logger     *zerolog.Logger
}

func LoadService(r AssignmentRepository, l *zerolog.Logger) *AssignmentService {
	return &AssignmentService{
		Logger:     l,
		Repository: r,
	}
}

func (as *AssignmentService) ApplyAssignment(request assignmentEntity.Assignment, personId uint) error {
	as.Logger.Info().Msg("Applying assignment")

	if err := request.Validate(); err != nil {
		return errors.New(err.Error())
	}

	assignment, err := as.Repository.GetAssignmentsByPersonAndFeatureFlagId(request.PersonID, request.FeatureFlagID)
	if err != nil {
		return err
	}

	if assignment.ID != 0 {
		return errors.New(fmt.Sprintf("Person %d is already assigned to the feature flag %d", request.PersonID, request.FeatureFlagID))
	}

	// TODO: validate ids against DB
	return as.Repository.ApplyAssignment(model.Assignment{
		PersonID:      request.PersonID,
		FeatureFlagID: request.FeatureFlagID,
	})
}

func (as *AssignmentService) DeleteAssignment(request assignmentEntity.Assignment, personId uint) error {
	as.Logger.Info().Msg("Delete assignment")

	if err := request.Validate(); err != nil {
		return errors.New(err.Error())
	}

	assignment, err := as.Repository.GetAssignmentsByPersonAndFeatureFlagId(request.PersonID, request.FeatureFlagID)
	if err != nil {
		return err
	}

	if assignment.ID == 0 {
		return errors.New(fmt.Sprintf("Person %d is not assigned to the feature flag %d", request.PersonID, request.FeatureFlagID))
	}

	return as.Repository.DeleteAssignment(model.Assignment{
		PersonID:      request.PersonID,
		FeatureFlagID: request.FeatureFlagID,
	})
}

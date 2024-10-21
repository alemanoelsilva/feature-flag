package mysql

import (
	"ff/internal/assignment"
	"ff/internal/db/repository"
	featureflag "ff/internal/feature_flag"
	"ff/internal/person"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewSqlFeatureFlagRepository(db *gorm.DB, logger *zerolog.Logger) featureflag.FeatureFlagRepository {
	featureFlagRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &featureFlagRepository
}

func NewSqlAssignmentRepository(db *gorm.DB, logger *zerolog.Logger) assignment.AssignmentRepository {
	assignmentRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &assignmentRepository
}

func NewSqlPersonRepository(db *gorm.DB, logger *zerolog.Logger) person.PersonRepository {
	personRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &personRepository
}

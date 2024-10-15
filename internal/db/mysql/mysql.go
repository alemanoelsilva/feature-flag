package mysql

import (
	"ff/internal/db/repository"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewSqlFeatureFlagRepository(db *gorm.DB, logger *zerolog.Logger) repository.FeatureFlagRepository {
	featureFlagRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &featureFlagRepository
}

func NewSqlAssignmentRepository(db *gorm.DB, logger *zerolog.Logger) repository.AssignmentRepository {
	assignmentRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &assignmentRepository
}

func NewSqlPersonRepository(db *gorm.DB, logger *zerolog.Logger) repository.PersonRepository {
	personRepository := repository.SqlRepository{DB: db, Logger: logger}
	return &personRepository
}

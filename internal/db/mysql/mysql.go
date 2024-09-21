package mysql

import (
	"ff/internal/db/repository"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewSqlRepository(db *gorm.DB, logger *zerolog.Logger) *repository.SqlRepository {
	featureFlagRepository := repository.SqlRepository{DB: db, Logger: logger}

	return &featureFlagRepository
}

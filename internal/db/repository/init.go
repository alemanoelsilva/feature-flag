package repository

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SqlRepository struct {
	DB     *gorm.DB
	Logger *zerolog.Logger
}

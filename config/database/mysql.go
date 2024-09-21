package database

import (
	"ff/internal/db/model"
	"os"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DDB struct {
	Logger *zerolog.Logger
}

func (ddb *DDB) Connect(uri string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(uri), &gorm.Config{})

	if err != nil {
		ddb.Logger.Fatal().Err(err).Msg("MySQL connection error")
		os.Exit(1)
	}

	return db
}

// TODO: Take a look at this
func (ddb *DDB) RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&model.FeatureFlag{}, &model.Person{})
}

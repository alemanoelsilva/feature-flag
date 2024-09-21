package main

import (
	"ff/config"
	"ff/config/database"
	"fmt"
	"log"
	"net/http"
	"os"

	handler "ff/api/handlers/http"
	"ff/internal/db/mysql"
	featureflag "ff/internal/feature-flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	config.LoadAppConfig(&logger)

	logger.Info().Msg("Initializing DB (MySQL)")
	logger.Info().Msg(config.AppConfig.ConnectionString)
	ddb := database.DDB{Logger: &logger}
	db := ddb.Connect(config.AppConfig.ConnectionString)
	ddb.RunMigrations(db)

	logger.Info().Msg("Initializing Repository (MySQL)")
	// TODO: split repositories
	featureFlagRepository := mysql.NewSqlRepository(db, &logger)

	logger.Info().Msg("Initializing Services/UseCases")
	featureFlagService := featureflag.LoadService(*featureFlagRepository, &logger)

	logger.Info().Msg("Initializing Handlers")
	router := handler.NewEchoHandler(*featureFlagService)

	// Start the server
	logger.Info().Msg(fmt.Sprintf("Starting Server on port %s", config.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.Port), router))
}

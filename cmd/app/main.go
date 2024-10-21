package main

import (
	"ff/config"
	"ff/config/database"
	"fmt"
	"log"
	"net/http"
	"os"

	handler "ff/api/handlers/http"
	assignment "ff/internal/assignment"
	mysql "ff/internal/db/mysql"
	featureflag "ff/internal/feature_flag"
	person "ff/internal/person"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	config.LoadAppConfig(&logger)
	ddb := database.DDB{Logger: &logger}

	logger.Info().Msg("Initializing DB (MySQL)")
	db := ddb.Connect(config.AppConfig.ConnectionString)
	ddb.RunMigrations(db)

	logger.Info().Msg("Initializing Repository (MySQL)")
	featureFlagRepository := mysql.NewSqlFeatureFlagRepository(db, &logger)
	assignmentRepository := mysql.NewSqlAssignmentRepository(db, &logger)
	peopleRepository := mysql.NewSqlPersonRepository(db, &logger)

	logger.Info().Msg("Initializing Services/UseCases")
	featureFlagService := featureflag.LoadService(featureFlagRepository, &logger)
	assignmentService := assignment.LoadService(assignmentRepository, &logger)
	personService := person.LoadService(peopleRepository, &logger)

	e := echo.New()
	e.Use(middleware.Logger())

	logger.Info().Msg("Initializing Handlers")
	handler.NewFeatureFlagEchoHandler(featureFlagService, e)
	handler.NewAssignmentEchoHandler(assignmentService, e)
	handler.NewPersonEchoHandler(personService, e)

	// Start the server
	logger.Info().Msg(fmt.Sprintf("Starting Server on port %s", config.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.Port), e))
}

package main

import (
	"ff/api/middlewares"
	"ff/config"
	"ff/config/database"
	assignment "ff/internal/assignment"
	"ff/internal/db/mysql"
	featureflag "ff/internal/feature_flag"
	person "ff/internal/person"
	handler "ff/web/handlers"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func loadServices() (*featureflag.FeatureFlagService, *assignment.AssignmentService, *person.PeopleService) {
	logger := zerolog.New(os.Stdout)

	config.LoadAppConfig(&logger)
	ddb := database.DDB{Logger: &logger}

	db := ddb.Connect(config.AppConfig.ConnectionString)
	ddb.RunMigrations(db)

	featureFlagRepository := mysql.NewSqlFeatureFlagRepository(db, &logger)
	assignmentRepository := mysql.NewSqlAssignmentRepository(db, &logger)
	peopleRepository := mysql.NewSqlPersonRepository(db, &logger)

	featureFlagService := featureflag.LoadService(featureFlagRepository, &logger)
	assignmentService := assignment.LoadService(assignmentRepository, &logger)
	personService := person.LoadService(peopleRepository, &logger)

	return featureFlagService, assignmentService, personService
}

const COOKIE_TEST = "HEEEEY FILL ME UP"

func TEST_AUTH(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie := &http.Cookie{
			Name:     "Cookie",
			Value:    COOKIE_TEST,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		}
		c.SetCookie(cookie)
		return next(c)
	}
}

func setupRoutes(e *echo.Echo) {
	featureFlagService, assignmentService, personService := loadServices()
	ffh := handler.FeatureFlagHandler{
		FeatureFlagService: featureFlagService,
	}
	ah := handler.AssignmentHandler{
		AssignmentService:  assignmentService,
		PersonService:      personService,
		FeatureFlagService: featureFlagService,
	}
	ch := handler.ComponentHandler{}

	e.GET("/", func(c echo.Context) error {
		// in this case, "/"  will be the same of "/feature-flags"
		return ffh.GetFeatureFlagList(c)
	}, TEST_AUTH)

	g := e.Group(("/feature-flags"), TEST_AUTH, middlewares.ValidateCookie)

	//! Pages
	g.GET("", ffh.GetFeatureFlagList)
	g.GET("/", ffh.GetFeatureFlagList)
	g.GET("/:id/assignments", ah.GetPeopleListToAssign)
	g.GET("/form/create-or-update", ffh.GetCreateOrUpdateFeatureFlag)
	// g.GET("/audit", adth.GetFeatureFlagList)

	//! Actions
	//* feature flag handlers
	g.POST("", ffh.CreateFeatureFlag)
	g.PUT("/:id", ffh.UpdateFeatureFlag)
	g.GET("/filters", ffh.GetFeatureFlagListFiltered)
	g.PUT("/status/:id", ffh.UpdateFeatureFlagStatus)

	//* assignment handlers
	g.GET("/:feature-flag-id/assignments/filters", ah.GetPeopleListToAssignFiltered)
	g.PUT("/:feature-flag-id/assignments/:id", ah.UpdateAssignment)
	g.PUT("/:feature-flag-id/global", ah.SetFeatureFlagToGlobal)

	//! Specific components updated by event
	//* is_global_event
	g.GET("/:feature-flag-id/component/set-global-button", ah.GetGlobalButtonSetup)
	g.GET("/:feature-flag-id/component/show-only-assigned-people", ah.GetShowOnlyAssignedPeopleFilter)

	//* create_feature_flag_event
	// g.GET("/component/header", ch.GetHeader)

	//! Components rendering
	//* error message TODO: search more about it
	g.GET("/component/error/dismiss", ch.DismissErrorMessage)
	// g.GET("/component/modal", ch.OpenModal)
	// http://localhost:6969/feature-flags/component/modal

}

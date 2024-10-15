package http

import (
	"errors"
	middlewares "ff/api/middlewares"
	person "ff/internal/person"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PeopleEchoHandler struct {
	PeopleService person.PeopleService
}

func NewPersonEchoHandler(person person.PeopleService, e *echo.Echo) {
	handler := &PeopleEchoHandler{
		PeopleService: person,
	}

	LoadPeopleRoutes(e, handler)
}

func LoadPeopleRoutes(e *echo.Echo, handler *PeopleEchoHandler) {
	group := e.Group("/api/feature-flags")

	group.GET("/v1/people", handler.getPersonHandler, middlewares.ValidateCookie)
	group.GET("/v1/people/feature-flags/:id", handler.getPersonWithAssignmentHandler, middlewares.ValidateCookie)
	group.GET("/v1/people/:id/assigned-feature-flags", handler.getAssignedFeatureFlagsByPersonIdHandler)
}

func (e *PeopleEchoHandler) getPersonHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	name := c.QueryParam("name")
	authInfo := c.Get("auth_info")

	fmt.Println(authInfo)

	if page <= 1 {
		page = 1 // Default page
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}

	people, totalCount, err := e.PeopleService.GetPeople(page, limit, name)
	if err != nil {
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	// TODO: check it again, it is terrible
	interfaceSlice := make([]interface{}, len(people))
	for i, v := range people {
		interfaceSlice[i] = v
	}

	return response.PaginationHandler(interfaceSlice, totalCount)
}

func (e *PeopleEchoHandler) getPersonWithAssignmentHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("feature flag id is not a number"))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	name := c.QueryParam("name")

	// TODO: check it again, it is terrible
	isAssignedStr := c.QueryParam("isAssigned")
	var isAssigned *bool
	if isAssignedStr == "true" {
		trueValue := true
		isAssigned = &trueValue
	} else if isAssignedStr == "false" {
		falseValue := false
		isAssigned = &falseValue
	} else if isAssignedStr != "" {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("invalid isAssigned value"))
	}

	if page <= 1 {
		page = 1 // Default page
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}

	people, totalCount, err := e.PeopleService.GetPeopleAssignmentByFeatureFlag(page, limit, uint(id), name, isAssigned)
	if err != nil {
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	// TODO: check it again, it is terrible
	interfaceSlice := make([]interface{}, len(people))
	for i, v := range people {
		interfaceSlice[i] = v
	}

	return response.PaginationHandler(interfaceSlice, totalCount)
}

func (e *PeopleEchoHandler) getAssignedFeatureFlagsByPersonIdHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("person id is not a number"))
	}

	featureFlags, err := e.PeopleService.GetAssignedFeatureFlagsByPersonId(uint(id))
	if err != nil {
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	// TODO: check it again, it is terrible
	interfaceSlice := make([]interface{}, len(featureFlags))
	for i, v := range featureFlags {
		interfaceSlice[i] = v
	}

	return response.PaginationHandler(interfaceSlice, int64(len(featureFlags)))
}

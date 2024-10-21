package http

import (
	"errors"
	middlewares "ff/api/middlewares"
	"ff/internal/db/model"
	p_entity "ff/internal/person/entity"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PersonService interface {
	GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, filters p_entity.PersonFilters) ([]p_entity.PersonWithAssignmentResponse, int64, error)
	GetAssignedFeatureFlagsByPersonId(id uint) ([]p_entity.AssignedFeatureFlagResponse, error)
}

type PeopleEchoHandler struct {
	PeopleService PersonService
}

func NewPersonEchoHandler(person PersonService, e *echo.Echo) {
	handler := &PeopleEchoHandler{
		PeopleService: person,
	}

	LoadPeopleRoutes(e, handler)
}

func LoadPeopleRoutes(e *echo.Echo, handler *PeopleEchoHandler) {
	group := e.Group("/api/feature-flags")

	group.GET("/v1/people/feature-flags/:id", handler.getPersonWithAssignmentHandler, middlewares.ValidateCookie)
	group.GET("/v1/people/:id/assigned-feature-flags", handler.getAssignedFeatureFlagsByPersonIdHandler)
	// TODO: get feature flag / 1/ person / 1/ to get a single register? make sense?
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

	pagination := model.Pagination{
		Page:  page,
		Limit: limit,
	}

	filters := p_entity.PersonFilters{
		FeatureFlagID: uint(id),
		Name:          name,
		IsAssigned:    isAssigned,
	}

	people, totalCount, err := e.PeopleService.GetPeopleAssignmentByFeatureFlag(pagination, filters)
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

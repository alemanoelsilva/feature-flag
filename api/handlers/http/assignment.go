package http

import (
	"ff/api/middlewares"
	a_entity "ff/internal/assignment/entity"
	"ff/pkg/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AssignmentService interface {
	ApplyAssignment(request a_entity.Assignment, personId uint) error
	DeleteAssignment(request a_entity.Assignment, personId uint) error
}

type AssignmentEchoHandler struct {
	AssignmentService AssignmentService
}

func NewAssignmentEchoHandler(assignment AssignmentService, e *echo.Echo) {
	handler := &AssignmentEchoHandler{
		AssignmentService: assignment,
	}

	LoadAssignmentRoutes(e, handler)
}

func LoadAssignmentRoutes(e *echo.Echo, handler *AssignmentEchoHandler) {
	group := e.Group("/api/feature-flags", middlewares.ValidateCookie)

	group.POST("/v1/assignments", handler.applyAssignmentsHandler)
	group.DELETE("/v1/assignments", handler.removeAssignmentsHandler)
}

func (e *AssignmentEchoHandler) applyAssignmentsHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input a_entity.Assignment
	if err := utils.GetBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	var personId int
	if err := utils.GetAuthenticatedPerson(c, &personId); err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.AssignmentService.ApplyAssignment(input, uint(personId)); err != nil {
		alreadyAssignedError := fmt.Sprintf("Person %d is already assigned to the feature flag %d", input.PersonID, input.FeatureFlagID)
		if err.Error() == alreadyAssignedError {
			return response.SuccessHandlerMessage(http.StatusConflict, alreadyAssignedError)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandlerMessage(http.StatusCreated, "Assignment Applied")
}

func (e *AssignmentEchoHandler) removeAssignmentsHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input a_entity.Assignment
	if err := utils.GetBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	var personId int
	if err := utils.GetAuthenticatedPerson(c, &personId); err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.AssignmentService.DeleteAssignment(input, uint(personId)); err != nil {
		featureFlagNotAssigned := fmt.Sprintf("Person %d is not assigned to the feature flag %d", input.PersonID, input.FeatureFlagID)
		if err.Error() == featureFlagNotAssigned {
			return response.SuccessHandlerMessage(http.StatusConflict, featureFlagNotAssigned)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandlerMessage(http.StatusCreated, "Assignment Removed")
}

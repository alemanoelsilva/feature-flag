package handler

import (
	"errors"
	a_entity "ff/internal/assignment/entity"
	"ff/internal/auth"
	"ff/internal/db/model"
	ff_entity "ff/internal/feature_flag/entity"
	p_entity "ff/internal/person/entity"
	"ff/web/components"
	"ff/web/utils"
	"ff/web/views"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AssignmentService interface {
	ApplyAssignment(request a_entity.Assignment, personId uint) error
	DeleteAssignment(request a_entity.Assignment, personId uint) error
}

type PersonService interface {
	GetPeopleAssignmentByFeatureFlag(pagination model.Pagination, filters p_entity.PersonFilters) ([]p_entity.PersonWithAssignmentResponse, int64, error)
	GetAssignedFeatureFlagsByPersonId(id uint) ([]p_entity.AssignedFeatureFlagResponse, error)
}

type AssignmentHandler struct {
	AssignmentService  AssignmentService
	PersonService      PersonService
	FeatureFlagService FeatureFlagService
}

func FindPersonByID(id int, assignments *[]p_entity.PersonWithAssignmentResponse) p_entity.PersonWithAssignmentResponse {
	var assignment p_entity.PersonWithAssignmentResponse
	for _, p := range *assignments {
		if p.ID == strconv.Itoa(id) {
			assignment = p
		}
	}

	return assignment
}

func (ah *AssignmentHandler) GetPeopleListToAssign(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Response().Header().Add("HX-Replace-Url", "/error")
		return utils.Render(c, http.StatusBadRequest, views.GenericErrorPage("Feature Flag id is invalid (not a number)"))
	}

	featureFlags, total, _ := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})
	if total == 0 {
		c.Response().Header().Add("HX-Replace-Url", "/404")
		return utils.Render(c, http.StatusNotFound, views.NotFoundPage("Feature Flag not found"))
	}

	assignments, total, err := ah.PersonService.GetPeopleAssignmentByFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 500,
	}, p_entity.PersonFilters{
		FeatureFlagID: uint(id),
	})
	if err != nil {
		c.Response().Header().Add("HX-Replace-Url", "/error")
		return utils.Render(c, http.StatusPreconditionFailed, views.GenericErrorPage("Something goes wrong when attempting to get the assignment list"))
	}
	if total == 0 {
		c.Response().Header().Add("HX-Replace-Url", "/404")
		return utils.Render(c, http.StatusNotFound, views.NotFoundPage("Feature Flag not found"))
	}

	return utils.Render(c, http.StatusOK, views.AssignmentsPage(assignments, featureFlags[0]))
}

func (ah *AssignmentHandler) GetPeopleListToAssignFiltered(c echo.Context) error {
	idStr := c.Param("feature-flag-id")
	fmt.Printf("id %v \n", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorMessage(c, "Feature flag id is invalid (not a number)")
	}

	name := c.QueryParams().Get("name")
	isAssignedStr := c.QueryParams().Get("isAssigned")

	filters := p_entity.PersonFilters{
		FeatureFlagID: uint(id),
		Name:          name,
	}

	if isAssignedStr == "on" {
		isAssigned := true
		filters.IsAssigned = &isAssigned
	}

	assignments, _, err := ah.PersonService.GetPeopleAssignmentByFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 500,
	}, filters)
	if err != nil {
		return utils.ErrorMessage(c, "Something goes wrong when attempting to get the assignment list")
	}

	featureFlags, _, _ := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})

	return utils.Render(c, http.StatusOK, components.AssignmentTable(assignments, featureFlags[0]))
}

func (ah *AssignmentHandler) UpdateAssignment(c echo.Context) error {
	personIdStr := c.Param("id")
	personId, err := strconv.Atoi(personIdStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	featureFlagIdStr := c.Param("feature-flag-id")
	featureFlagId, err := strconv.Atoi(featureFlagIdStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	assignments, total, err := ah.PersonService.GetPeopleAssignmentByFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 10000,
	}, p_entity.PersonFilters{
		FeatureFlagID: uint(featureFlagId),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the assignment list")
	}
	if total == 0 {
		return errors.New("Feature Flag ID is invalid")
	}

	personAssignment := FindPersonByID(personId, &assignments)

	if personAssignment.IsAssigned {
		ah.AssignmentService.DeleteAssignment(a_entity.Assignment{
			PersonID:      uint(personId),
			FeatureFlagID: uint(featureFlagId),
		}, uint(authInfo.PersonID))
	} else {
		ah.AssignmentService.ApplyAssignment(a_entity.Assignment{
			PersonID:      uint(personId),
			FeatureFlagID: uint(featureFlagId),
		}, uint(authInfo.PersonID))
	}

	personAssignment.IsAssigned = !personAssignment.IsAssigned

	featureFlags, _, err := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(featureFlagId),
	})

	name := c.FormValue("name")
	isAssignedStr := c.FormValue("isAssigned")

	filters := p_entity.PersonFilters{
		FeatureFlagID: uint(featureFlagId),
		Name:          name,
	}

	if isAssignedStr == "on" {
		isAssigned := true
		filters.IsAssigned = &isAssigned
	}

	assignmentsToShow, _, _ := ah.PersonService.GetPeopleAssignmentByFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 10000,
	}, filters)

	return utils.Render(c, http.StatusOK, components.AssignmentTable(assignmentsToShow, featureFlags[0]))
}

func (ah *AssignmentHandler) SetFeatureFlagToGlobal(c echo.Context) error {
	featureFlagIdStr := c.Param("feature-flag-id")
	featureFlagId, err := strconv.Atoi(featureFlagIdStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	// authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	featureFlags, total, err := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(featureFlagId),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the assignment list")
	}
	if total == 0 {
		return errors.New("Feature Flag ID is invalid")
	}

	if err := ah.FeatureFlagService.UpdateFeatureFlagById(uint(featureFlagId), ff_entity.UpdateFeatureFlag{
		Description:    featureFlags[0].Description,
		IsActive:       featureFlags[0].IsActive,
		IsGlobal:       !featureFlags[0].IsGlobal,
		ExpirationDate: featureFlags[0].ExpirationDate,
	}); err != nil {
		return errors.New("Something goes wrong when attempting to update the feature flag global")
	}

	featureFlags[0].IsGlobal = !featureFlags[0].IsGlobal

	name := c.FormValue("name")
	isAssignedStr := c.FormValue("isAssigned")

	filters := p_entity.PersonFilters{
		FeatureFlagID: uint(featureFlagId),
		Name:          name,
	}

	if isAssignedStr == "on" {
		isAssigned := true
		filters.IsAssigned = &isAssigned
	}

	assignmentsToShow, _, _ := ah.PersonService.GetPeopleAssignmentByFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 10000,
	}, filters)

	// c.Response().Header().Add("HX-Trigger-After-Swap", `{"isGlobal":{"target":"#is_global_button"}}`)
	c.Response().Header().Add("HX-Trigger-After-Swap", "is_global_event")

	return utils.Render(c, http.StatusOK, components.AssignmentTable(assignmentsToShow, featureFlags[0]))
}

func (ah *AssignmentHandler) GetGlobalButtonSetup(c echo.Context) error {
	fmt.Print("\n\nGEEEET SET GLOBAL BUTTON\n\n")

	featureFlagIdStr := c.Param("feature-flag-id")
	featureFlagId, err := strconv.Atoi(featureFlagIdStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	// authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	featureFlags, total, err := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(featureFlagId),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the feature flag list")
	}
	if total == 0 {
		return errors.New("Feature Flag ID is invalid")
	}

	// c.Response().Header().Add("HX-Trigger", "isGlobal")

	return utils.Render(c, http.StatusOK, components.IsGlobalButton(featureFlags[0]))
}

func (ah *AssignmentHandler) GetShowOnlyAssignedPeopleFilter(c echo.Context) error {
	fmt.Print("\n\nGEEEET ONLY ASSIGNED PEOPLE CHECKBOX\n\n")

	featureFlagIdStr := c.Param("feature-flag-id")
	featureFlagId, err := strconv.Atoi(featureFlagIdStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	// authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	featureFlags, total, err := ah.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(featureFlagId),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the feature flag list")
	}
	if total == 0 {
		return errors.New("Feature Flag ID is invalid")
	}

	// c.Response().Header().Add("HX-Trigger", "isGlobal")

	return utils.Render(c, http.StatusOK, components.ShowOnlyAssignedPeopleFilter(featureFlags[0]))
}

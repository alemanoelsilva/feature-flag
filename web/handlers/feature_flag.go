package handler

import (
	"ff/internal/auth"
	"ff/internal/db/model"
	ff_entity "ff/internal/feature_flag/entity"
	"ff/web/components"
	"ff/web/types"
	"ff/web/utils"
	"ff/web/views"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func FindFeatureFlagByID(id int, ffr *[]ff_entity.FeatureFlagResponse) ff_entity.FeatureFlagResponse {
	var featureFlags ff_entity.FeatureFlagResponse
	for _, ff := range *ffr {
		intID, _ := strconv.Atoi(ff.ID)
		if intID == id {
			featureFlags = ff
		}
	}

	return featureFlags
}

type FeatureFlagService interface {
	CreateFeatureFlag(request ff_entity.FeatureFlag, personId uint) error
	GetFeatureFlag(pagination model.Pagination, filters ff_entity.FeatureFlagFilters) ([]ff_entity.FeatureFlagResponse, int64, error)
	UpdateFeatureFlagById(id uint, request ff_entity.UpdateFeatureFlag) error
}

type FeatureFlagHandler struct {
	FeatureFlagService FeatureFlagService
}

func (ffh *FeatureFlagHandler) GetFeatureFlagList(c echo.Context) error {
	// ff := services.FeatureFlag{}

	// TODO: user ffh.Service

	featureFlags, _, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, ff_entity.FeatureFlagFilters{})

	if err != nil {
		return utils.ErrorMessage(c, "something goes wrong when attempting to get the feature flag list")
	}

	// featureFlags := ff.GetFeatureFlag()
	return utils.Render(c, http.StatusOK, views.FeatureFlagsPage(featureFlags))
}

func (ffh *FeatureFlagHandler) GetCreateOrUpdateFeatureFlag(c echo.Context) error {
	idStr := c.QueryParam("id")

	// return create form
	if idStr == "" {
		// c.Response().Header().Add("HX-Trigger-After-Swap", "create_feature_flag_event")
		return utils.Render(c, http.StatusOK, components.Modal(true, ff_entity.FeatureFlagResponse{}))
	}

	// otherwise, get feature flag and return the update form filled up
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Response().Header().Add("HX-Replace-Url", "/404")
		return utils.Render(c, http.StatusNotFound, views.NotFoundPage("feature flag ID is not a valid number"))
	}

	featureFlags, _, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})
	if err != nil {
		c.Response().Header().Add("HX-Replace-Url", "/error")
		return utils.Render(c, http.StatusNotFound, views.NotFoundPage("something goes wrong when attempting to get the feature flag by id"))
	}

	// c.Response().Header().Add("HX-Trigger-After-Swap", "create_feature_flag_event")
	return utils.Render(c, http.StatusOK, components.Modal(true, featureFlags[0]))
}

func (ffh *FeatureFlagHandler) GetFeatureFlagListFiltered(c echo.Context) error {
	name := c.QueryParams().Get("name")
	isActiveStr := c.QueryParams().Get("isActive")

	// TODO: user ffh.Service

	filters := ff_entity.FeatureFlagFilters{
		Name: name,
	}

	if isActiveStr == "on" {
		isActive := true
		filters.IsActive = &isActive
	}

	// fmt.Printf("filters %v | %v", filters.Name, *filters.IsActive)
	featureFlags, _, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, filters)

	if err != nil {
		// return errors.New("something goes wrong when attempting to get the feature flag list")
		return utils.ErrorMessage(c, "something goes wrong when attempting to get the feature flag list")
	}

	return utils.Render(c, http.StatusOK, components.FeatureFlagTable(featureFlags))
}

func (ffh *FeatureFlagHandler) UpdateFeatureFlagStatus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorMessage(c, "feature flag id is invalid (not a number)")
	}

	featureFlags, total, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})
	if err != nil {
		return utils.ErrorMessage(c, "something goes wrong when attempting to get the feature flag list")
	}
	if total == 0 {
		return utils.ErrorMessage(c, "feature Flag ID is invalid")
	}

	selectedFeatureFlag := FindFeatureFlagByID(id, &featureFlags)
	selectedFeatureFlag.IsActive = !selectedFeatureFlag.IsActive

	// TODO: check if selectedFeatureFlag exists
	requestToUpdate := ff_entity.UpdateFeatureFlag{
		Description:    selectedFeatureFlag.Description,
		IsActive:       selectedFeatureFlag.IsActive,
		IsGlobal:       selectedFeatureFlag.IsGlobal,
		ExpirationDate: selectedFeatureFlag.ExpirationDate,
	}

	err = ffh.FeatureFlagService.UpdateFeatureFlagById(uint(id), requestToUpdate)
	if err != nil {
		return utils.ErrorMessage(c, "something goes wrong when attempting to update the feature flag")
	}

	name := c.FormValue("name")
	isActiveStr := c.FormValue("isActive")

	filters := ff_entity.FeatureFlagFilters{
		Name: name,
	}

	if isActiveStr == "on" {
		isActive := true
		filters.IsActive = &isActive
	}

	response, _, _ := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, filters)

	return utils.Render(c, http.StatusOK, components.FeatureFlagTable(response))
}

func (ffh *FeatureFlagHandler) CreateFeatureFlag(c echo.Context) error {
	name := strings.Trim(c.FormValue("name"), " ")
	description := strings.Trim(c.FormValue("description"), " ")
	isActive := c.FormValue("isActive") == "on"
	expirationDate := c.FormValue("expirationDate")

	authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	err := ffh.FeatureFlagService.CreateFeatureFlag(ff_entity.FeatureFlag{
		Name:           name,
		Description:    description,
		IsActive:       isActive,
		ExpirationDate: expirationDate,
	}, uint(authInfo.PersonID))

	// error on feature flag creation
	if err != nil {
		errorType := strings.Split(err.Error(), "|")[0]

		var errorResponse types.ErrorCreateFeatureFlagForm

		if errorType == "Name" {
			errorResponse.HasError = true
			errorResponse.IsNameError = true
			errorResponse.ErrorMessage = strings.Split(err.Error(), "|")[1]
			c.Response().Header().Add("HX-Trigger", "isNameErrorEvent")
		} else if errorType == "Description" {
			errorResponse.HasError = true
			errorResponse.IsDescriptionError = true
			errorResponse.ErrorMessage = strings.Split(err.Error(), "|")[1]
			c.Response().Header().Add("HX-Trigger", "isDescriptionErrorEvent")
		} else {
			errorResponse.HasError = true
			errorResponse.IsRequestError = true
			errorResponse.ErrorMessage = err.Error()
		}

		return utils.ErrorMessage(c, errorResponse.ErrorMessage)
	}

	c.Response().Header().Add("HX-Retarget", "#message")
	c.Response().Header().Add("HX-Trigger", "closeModal")
	c.Response().Header().Add("HX-Trigger", "refresh_ff_list_event")
	return utils.Render(c, http.StatusConflict, components.Message(true, "Feature Flag created", false))
}

func (ffh *FeatureFlagHandler) UpdateFeatureFlag(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorMessage(c, "feature flag ID is not a valid number")
	}

	name := c.FormValue("name")
	description := strings.Trim(c.FormValue("description"), " ")
	isActive := c.FormValue("isActive") == "on"
	expirationDate := c.FormValue("expirationDate")

	ffOnDB, _, _ := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		Name: name,
	})

	err = ffh.FeatureFlagService.UpdateFeatureFlagById(uint(id), ff_entity.UpdateFeatureFlag{
		// method updates all 4 fields, getting the current isGlobal value to not set false when it is true
		Description:    description,
		IsActive:       isActive,
		IsGlobal:       ffOnDB[0].IsGlobal,
		ExpirationDate: expirationDate,
	})

	// error on feature flag creation
	if err != nil && err.Error() != "no feature flag updated" {
		ff := ffOnDB[0]

		ff.Description = description
		ff.IsActive = isActive
		ff.ExpirationDate = expirationDate

		errorType := strings.Split(err.Error(), "|")[0]
		// errorMessage := strings.Split(err.Error(), "|")[1]

		var errorResponse types.ErrorCreateFeatureFlagForm

		if errorType == "Description" {
			errorResponse.HasError = true
			errorResponse.IsDescriptionError = true
			errorResponse.ErrorMessage = strings.Split(err.Error(), "|")[1]
			c.Response().Header().Add("HX-Trigger", "isDescriptionErrorEvent")
		} else {
			errorResponse.HasError = true
			errorResponse.IsRequestError = true
			errorResponse.ErrorMessage = err.Error()
		}

		return utils.ErrorMessage(c, errorResponse.ErrorMessage)
	}

	c.Response().Header().Add("HX-Retarget", "#message")
	c.Response().Header().Add("HX-Trigger", "closeModal") // Trigger closing the modal
	c.Response().Header().Add("HX-Trigger", "refresh_ff_list_event")
	return utils.Render(c, http.StatusConflict, components.Message(true, "Feature Flag updated", false))
}

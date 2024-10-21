package handler

import (
	"errors"
	"ff/internal/auth"
	"ff/internal/db/model"
	ff_entity "ff/internal/feature_flag/entity"
	"ff/web/components"
	"ff/web/types"
	"ff/web/utils"
	"ff/web/views"
	"fmt"
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
		return errors.New("Something goes wrong when attempting to get the feature flag list")
	}

	// featureFlags := ff.GetFeatureFlag()
	return utils.Render(c, http.StatusOK, views.FeatureFlagsPage(featureFlags))
}

func (ffh *FeatureFlagHandler) GetCreateOrUpdateFeatureFlag(c echo.Context) error {
	idStr := c.QueryParam("id")

	// return create form
	if idStr == "" {
		c.Response().Header().Add("HX-Trigger-After-Swap", "create_feature_flag_event")

		return utils.Render(c, http.StatusOK, views.CreateOrUpdateFeatureFlagPage(ff_entity.FeatureFlagResponse{}, types.ErrorCreateFeatureFlagForm{}))
	}

	// otherwise, get feature flag and return the update form filled up
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("Feature flag ID is not a valid number")
	}

	featureFlags, _, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 1,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the feature flag by id")
	}

	c.Response().Header().Add("HX-Trigger-After-Swap", "create_feature_flag_event")

	return utils.Render(c, http.StatusOK, views.CreateOrUpdateFeatureFlagPage(featureFlags[0], types.ErrorCreateFeatureFlagForm{}))
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
		return errors.New("Something goes wrong when attempting to get the feature flag list")
	}

	return utils.Render(c, http.StatusOK, components.FeatureFlagTable(featureFlags))
}

func (ffh *FeatureFlagHandler) UpdateFeatureFlagStatus(c echo.Context) error {
	idStr := c.Param("id")
	fmt.Printf("id %v \n", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("Feature flag id is invalid (not a number)")
	}

	featureFlags, total, err := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, ff_entity.FeatureFlagFilters{
		ID: uint(id),
	})
	if err != nil {
		return errors.New("Something goes wrong when attempting to get the feature flag list")
	}
	if total == 0 {
		return errors.New("Feature Flag ID is invalid")
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
		return errors.New("Something goes wrong when attempting to update the feature flag")
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
	name := c.FormValue("name")
	description := c.FormValue("description")
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
		ff := ff_entity.FeatureFlagResponse{
			Name:           name,
			Description:    description,
			IsActive:       isActive,
			ExpirationDate: expirationDate,
		}

		errorType := strings.Split(err.Error(), "|")[0]
		errorMessage := strings.Split(err.Error(), "|")[1]

		var resultError types.ErrorCreateFeatureFlagForm

		if errorType == "Name" {
			resultError.HasError = true
			resultError.IsNameError = true
			resultError.ErrorMessage = errorMessage
		} else if errorType == "Description" {
			resultError.HasError = true
			resultError.IsDescriptionError = true
			resultError.ErrorMessage = errorMessage
		} else {
			resultError.HasError = true
			resultError.IsRequestError = true
			resultError.ErrorMessage = err.Error()
		}

		c.Response().Header().Add("HX-Retarget", "#create_feature_flag_page")
		return utils.Render(c, http.StatusConflict, views.CreateOrUpdateFeatureFlagPage(ff, resultError))
	}

	featureFlags, _, _ := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, ff_entity.FeatureFlagFilters{})

	return utils.Render(c, http.StatusOK, views.FeatureFlagsPage(featureFlags))
}

func (ffh *FeatureFlagHandler) UpdateFeatureFlag(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("Feature flag ID is not a valid number")
	}

	name := c.FormValue("name")
	description := c.FormValue("description")
	isActive := c.FormValue("isActive") == "on"
	expirationDate := c.FormValue("expirationDate")

	err = ffh.FeatureFlagService.UpdateFeatureFlagById(uint(id), ff_entity.UpdateFeatureFlag{
		Description:    description,
		IsActive:       isActive,
		ExpirationDate: expirationDate,
	})

	// error on feature flag creation
	if err != nil {
		ff := ff_entity.FeatureFlagResponse{
			Name:           name,
			Description:    description,
			IsActive:       isActive,
			ExpirationDate: expirationDate,
		}

		errorType := strings.Split(err.Error(), "|")[0]
		errorMessage := strings.Split(err.Error(), "|")[1]

		var resultError types.ErrorCreateFeatureFlagForm

		if errorType == "Description" {
			resultError.HasError = true
			resultError.IsDescriptionError = true
			resultError.ErrorMessage = errorMessage
		} else {
			resultError.HasError = true
			resultError.IsRequestError = true
			resultError.ErrorMessage = err.Error()
		}

		c.Response().Header().Add("HX-Retarget", "#create_feature_flag_page")
		return utils.Render(c, http.StatusConflict, views.CreateOrUpdateFeatureFlagPage(ff, resultError))
	}

	featureFlags, _, _ := ffh.FeatureFlagService.GetFeatureFlag(model.Pagination{
		Page:  1,
		Limit: 100,
	}, ff_entity.FeatureFlagFilters{})

	return utils.Render(c, http.StatusOK, views.FeatureFlagsPage(featureFlags))
}

package handler

import (
	"ff/web/components"
	"ff/web/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ComponentHandler struct {
}

// func (ch *ComponentHandler) GetHeader(c echo.Context) error {
// 	showCreateButton := c.QueryParam("showCreateButton") == "true"

// 	return utils.Render(c, http.StatusOK, components.CreateFeatureFlagButton(showCreateButton))
// }

func (ch *ComponentHandler) DismissErrorMessage(c echo.Context) error {
	return utils.Render(c, http.StatusOK, components.Message(false, "", false))
}

// func (ch *ComponentHandler) OpenModal(c echo.Context) error {
// 	return utils.Render(c, http.StatusOK, components.Modal(true))
// }

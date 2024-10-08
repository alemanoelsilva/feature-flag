package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseJSON struct {
	c echo.Context
}

func (s ResponseJSON) SuccessHandler(code int, data interface{}) error {
	return s.c.JSON(code, data)
}

func (s ResponseJSON) SuccessHandlerMessage(code int, msg string) error {
	return s.c.JSON(code, map[string]interface{}{"message": msg})
}

type PaginationResponse struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total"`
}

func (s ResponseJSON) PaginationHandler(data []interface{}, totalCount int64) error {
	return s.c.JSON(http.StatusOK, PaginationResponse{
		Items: data,
		Total: int(totalCount),
	})
}

func (s ResponseJSON) ErrorHandler(code int, err error) error {
	return s.c.JSON(code, map[string]interface{}{"error": err.Error()})
}

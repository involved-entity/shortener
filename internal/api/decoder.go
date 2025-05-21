package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func DecodeRequest(c echo.Context, dto interface{}) error {
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Response{Msg: "invalid data"})
	}
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Response{Msg: "invalid data"})
	}
	return nil
}

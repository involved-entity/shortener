package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DecodeRequest(c echo.Context, dto interface{}) error {
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: "invalid data"})
	}
	return nil
}

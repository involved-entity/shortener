package api

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func OptionalJWT(skippedURLs []string) func(c echo.Context) bool {
	return func(c echo.Context) bool {
		path := c.Request().URL.Path
		for _, url := range skippedURLs {
			if strings.HasPrefix(path, url) && len(c.Request().Header["Authorization"]) == 0 {
				return true
			}
		}

		return false
	}
}

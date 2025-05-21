package api

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func OptionalJWT(skippedURLs []string) func(c echo.Context) bool {
	return func(c echo.Context) bool {
		path := c.Path()
		for _, url := range skippedURLs {
			if strings.HasPrefix(url, path) {
				return true
			}
		}

		return false
	}
}

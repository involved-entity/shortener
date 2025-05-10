package urls

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
)

func DecodeClickRequest(c echo.Context) (string, string, string, string) {
	shortCode := c.Param("shortCode")

	lang := c.Request().Header.Get("Accept-Language")
	ua := c.Request().UserAgent()
	referer := c.Request().Header.Get("Referer")

	langParts := strings.Split(lang, ";")
	langCode := langParts[0]

	uaParts := strings.Split(ua, " ")
	browser := ""
	for _, part := range uaParts {
		if strings.Contains(part, "Chrome") || strings.Contains(part, "Firefox") || strings.Contains(part, "Safari") || strings.Contains(part, "Edge") {
			browser = part
			break
		}
	}

	return shortCode, referer, langCode, browser
}

func GetUserID(c echo.Context) int {
	return int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"].(float64))
}

package urls

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
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

func isDigitsOnly(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func stringToIntStrict(s string) int {
	if s == "" || !isDigitsOnly(s) {
		return 1
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return num
}

func GetPage(c echo.Context) int {
	return stringToIntStrict(c.QueryParam("page"))
}

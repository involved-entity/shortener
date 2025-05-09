package urls

import (
	"shortener/internal/api"
	"shortener/internal/database"

	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"net/http"
)

type URLDTO struct {
	OriginalURL string `json:"originalURL"`
	ShortCode   string `json:"shortCode"`
}

func SaveURL(c echo.Context) error {
	var userID int

	if c.Get("user") != nil {
		idValue := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"]
		userID = int(idValue.(float64))
	}

	dto := URLDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}
	db := database.GetDB()
	r := Repository{db: db, UserId: userID}
	url, err := r.SaveURL(dto.OriginalURL, dto.ShortCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "shortcode already used"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: url})
}

func GetURL(c echo.Context) error {
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

	db := database.GetDB()
	r := Repository{db: db}
	url, id, err := r.GetURL(shortCode)
	if err != nil {
		return c.String(http.StatusBadRequest, "shortcode is not defined")
	}
	r.RegisterClick(id, c.RealIP(), referer, langCode, browser)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

func DeleteURL(c echo.Context) error {
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"].(float64))
	shortCode := c.Param("shortCode")
	db := database.GetDB()
	r := Repository{db: db, UserId: userID}
	if err := r.DeleteURL(shortCode); err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "shortcode is not defined"})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetMyURLs(c echo.Context) error {
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"].(float64))
	db := database.GetDB()
	r := Repository{db: db, UserId: userID}
	urls, err := r.GetUserURLs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: urls})
}

func GetURLClicks(c echo.Context) error {
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"].(float64))
	shortCode := c.Param("shortCode")
	db := database.GetDB()
	r := Repository{db: db, UserId: userID}
	clicks, err := r.GetURLClicks(shortCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: clicks})
}

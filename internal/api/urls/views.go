package urls

import (
	"shortener/internal/api"
	"shortener/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"net/http"
)

type URLDTO struct {
	OriginalURL string `json:"originalURL" validate:"required,url"`
	ShortCode   string `json:"shortCode" validate:"required,min=2,max=32"`
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
	r := Repository{db: database.GetDB(), UserId: userID}
	url, err := r.SaveURL(dto.OriginalURL, dto.ShortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "shortcode already used"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: url})
}

func GetURL(c echo.Context) error {
	shortCode, referer, langCode, browser := DecodeClickRequest(c)

	r := Repository{db: database.GetDB()}
	url, id, err := r.GetURL(shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "shortcode is not defined"})
	}
	r.RegisterClick(id, c.RealIP(), referer, langCode, browser)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

func DeleteURL(c echo.Context) error {
	userID := GetUserID(c)
	shortCode := c.Param("shortCode")
	r := Repository{db: database.GetDB(), UserId: userID}
	if err := r.DeleteURL(shortCode); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "shortcode is not defined"})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetMyURLs(c echo.Context) error {
	userID := GetUserID(c)
	page := GetPage(c)
	r := Repository{db: database.GetDB(), UserId: userID, Page: page}
	urls, err := r.GetUserURLs()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: urls})
}

func GetURLClicks(c echo.Context) error {
	userID := GetUserID(c)
	shortCode := c.Param("shortCode")
	page := GetPage(c)
	r := Repository{db: database.GetDB(), UserId: userID, Page: page}
	clicks, err := r.GetURLClicks(shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: clicks})
}

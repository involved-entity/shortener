package urls

import (
	"shortener/internal/api"
	"shortener/internal/database"

	"github.com/labstack/echo/v4"

	"net/http"
)

type URLDTO struct {
	OriginalURL string `json:"originalURL"`
	ShortCode   string `json:"shortCode"`
}

func SaveURL(c echo.Context) error {
	dto := URLDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}
	db := database.GetDB()
	url, err := database.SaveURL(db, dto.OriginalURL, dto.ShortCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "shortcode already used"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: url})
}

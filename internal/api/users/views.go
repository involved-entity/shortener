package users

import (
	"net/http"
	"shortener/internal/api"
	"shortener/internal/database"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	dto := UserDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "Cant hash this password"})
	}

	db := database.GetDB()
	r := Repository{db: db}
	user, err := r.SaveUser(dto.Username, dto.Email, string(hashedPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "Username or email is already exists"})
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}

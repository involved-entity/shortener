package users

import (
	"net/http"
	"shortener/internal/api"
	"shortener/internal/database"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDTO struct {
	Token string `json:"token"`
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

func Login(ttl int, secret string) func(c echo.Context) error {
	return func(c echo.Context) error {
		dto := UserLoginDTO{}
		if err := api.DecodeRequest(c, &dto); err != nil {
			return err
		}

		db := database.GetDB()
		r := Repository{db: db}
		user, err := r.GetUser(dto.Username)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, api.Response{Msg: "User is not found or not verified"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid credentials")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(time.Minute * time.Duration(ttl)).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Could not generate token")
		}

		return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: TokenDTO{Token: tokenString}})
	}
}

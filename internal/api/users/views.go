package users

import (
	"log"
	"net/http"
	api "shortener/internal/api"
	conf "shortener/internal/config"
	"shortener/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type JWTData struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegenerateCodeDTO struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type VerificationDTO struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

type ResetPasswordDTO struct {
	Username string `json:"username"`
}

type ResetPasswordConfirmDTO struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	dto := UserDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	hashedPassword, err := GetHashedPassword(c, dto.Password)
	if err != nil {
		return err
	}

	db := database.GetDB()
	r := Repository{db: db}
	user, err := r.SaveUser(dto.Username, dto.Email, string(hashedPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "Username or email is already exists"})
	}

	if err := CreateAndSendToken(c, user.ID, user.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}

func Login(c echo.Context) error {
	dto := UserLoginDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	db := database.GetDB()
	r := Repository{db: db}
	user, err := r.GetUser(UserInfo{Username: dto.Username})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Response{Msg: "User is not found or not verified"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	config := conf.GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JWTData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		"exp": time.Now().Add(time.Minute * time.Duration(config.JWT.JWT_TTL)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT.SECRET))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate token")
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: TokenDTO{Token: tokenString}})
}

func RegenerateCode(c echo.Context) error {
	dto := RegenerateCodeDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}
	if err := CreateAndSendToken(c, uint(dto.ID), dto.Email); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success"})
}

func ActivateAccount(c echo.Context) error {
	dto := VerificationDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	config := conf.GetConfig()
	if err := CheckRedisToken(c, dto.ID, dto.Code, config.OTP.RedisName); err != nil {
		return err
	}

	db := database.GetDB()
	r := Repository{db: db}
	if err := r.VerificateUser(dto.ID); err != nil {
		log.Println("Error with database", err)
		return c.JSON(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}

	return c.JSON(http.StatusAccepted, api.Response{Msg: "success"})
}

func ResetPassword(c echo.Context) error {
	dto := ResetPasswordDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	db := database.GetDB()
	r := Repository{db: db}
	user, err := r.GetUser(UserInfo{Username: dto.Username})
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Msg: "User not found"})
	}

	if err := CreateAndSendResetPasswordLink(c, user.ID, user.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, api.Response{Msg: "success"})
}

func ResetPasswordConfirm(c echo.Context) error {
	dto := ResetPasswordConfirmDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	config := conf.GetConfig()
	if err := CheckRedisToken(c, dto.ID, dto.Token, config.ResetToken.RedisName); err != nil {
		return err
	}

	hashedPassword, err := GetHashedPassword(c, dto.Password)
	if err != nil {
		return err
	}

	db := database.GetDB()
	r := Repository{db: db}
	if err := r.ChangeUserPassword(dto.ID, string(hashedPassword)); err != nil {
		return c.JSON(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success"})
}

func GetMe(c echo.Context) error {
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(map[string]interface{})["id"].(float64))
	db := database.GetDB()
	r := Repository{db: db}
	user, err := r.GetUser(UserInfo{ID: userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Response{Msg: "User not found"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}

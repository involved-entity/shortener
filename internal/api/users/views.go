package users

import (
	"context"
	"log"
	"net/http"
	api "shortener/internal/api"
	"shortener/internal/database"
	"shortener/internal/redis"
	"strconv"
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

	if err := CreateAndSendToken(c, user.ID, user.Email); err != nil {
		return err
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
			"sub": JWTData{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
			"exp": time.Now().Add(time.Minute * time.Duration(ttl)).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Could not generate token")
		}

		return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: TokenDTO{Token: tokenString}})
	}
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

	redisClient := redis.GetClient()
	otp, err := redisClient.Get(context.Background(), "otp:"+strconv.Itoa(dto.ID)).Result()
	if err != nil {
		log.Println(otp)
		log.Printf("Error with redis: %T", err)
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "Code expired"})
	}

	if otp != dto.Code {
		return c.JSON(http.StatusBadRequest, api.Response{Msg: "Invalid code"})
	}

	db := database.GetDB()
	r := Repository{db: db}
	if err := r.VerificateUser(dto.ID); err != nil {
		log.Println("Error with database", err)
		return c.JSON(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}

	return c.JSON(http.StatusAccepted, api.Response{Msg: "success"})
}

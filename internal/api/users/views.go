package users

import (
	"context"
	"log"
	"net/http"
	api "shortener/internal/api"
	"shortener/internal/database"
	"shortener/internal/machinery"
	"shortener/internal/redis"
	"strconv"
	"time"

	machineryTasks "github.com/RichardKnop/machinery/v2/tasks"
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

	tokenOTP, err := GenerateSecureToken()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			api.Response{Msg: "Failed to generate token to verification. Please try again"},
		)
	}
	redisClient := redis.GetClient()
	redisClient.Set(context.Background(), "otp:"+strconv.Itoa(int(user.ID)), tokenOTP, time.Minute*5)

	machineryServer := machinery.GetServer()
	signature := &machineryTasks.Signature{
		Name: "send_email",
		Args: []machineryTasks.Arg{
			{Name: "email", Type: "string", Value: user.Email},
			{Name: "code", Type: "string", Value: tokenOTP},
		},
	}
	machineryServer.SendTaskWithContext(context.Background(), signature)

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

func ActivateAccount(c echo.Context) error {
	dto := VerificationDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	redisClient := redis.GetClient()
	otp, err := redisClient.Get(context.Background(), "code:"+strconv.Itoa(dto.ID)).Result()
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

package users

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net/http"
	"net/url"
	api "shortener/internal/api"
	conf "shortener/internal/config"
	"shortener/internal/machinery"
	"shortener/internal/redis"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	machineryTasks "github.com/RichardKnop/machinery/v2/tasks"
	"github.com/labstack/echo/v4"
)

func generateSecureToken(elements string, length int) (string, error) {
	token := make([]byte, length)

	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(elements))))
		if err != nil {
			log.Printf("failed to generate token: %v", err)
			return "", errors.New("failed to generate token")
		}
		token[i] = elements[num.Int64()]
	}

	return string(token), nil
}

func CreateAndSendToken(id uint, email string) error {
	tokenOTP, err := generateSecureToken("0123456789", 5)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			api.Response{Msg: "Failed to generate token to verification. Please try again"},
		)
	}
	redisClient := redis.GetClient()
	config := conf.GetConfig()
	redisClient.Set(
		context.Background(),
		config.OTP.RedisName+":"+strconv.Itoa(int(id)),
		tokenOTP,
		time.Minute*time.Duration(config.OTP.OTP_TTL),
	)

	machineryServer := machinery.GetServer()
	signature := &machineryTasks.Signature{
		Name: "send_email",
		Args: []machineryTasks.Arg{
			{Name: "email", Type: "string", Value: email},
			{Name: "code", Type: "string", Value: tokenOTP},
		},
	}
	machineryServer.SendTaskWithContext(context.Background(), signature)

	return nil
}

func CreateAndSendResetPasswordLink(id uint, email string) error {
	token, err := generateSecureToken("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			api.Response{Msg: "Failed to generate token to reset password. Please try again"},
		)
	}

	redisClient := redis.GetClient()
	config := conf.GetConfig()
	redisClient.Set(
		context.Background(),
		config.ResetToken.RedisName+":"+strconv.Itoa(int(id)),
		token,
		time.Minute*time.Duration(config.ResetToken.RT_TTL),
	)

	baseURL, err := url.Parse(config.ResetToken.FrontendUrl)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			api.Response{Msg: "Internal server error. Please try again"},
		)
	}

	query := url.Values{
		"token": {token},
		"id":    {strconv.Itoa(int(id))},
	}

	baseURL.RawQuery = query.Encode()

	machineryServer := machinery.GetServer()
	signature := &machineryTasks.Signature{
		Name: "reset_password",
		Args: []machineryTasks.Arg{
			{Name: "email", Type: "string", Value: email},
			{Name: "link", Type: "string", Value: baseURL.String()},
		},
	}
	machineryServer.SendTaskWithContext(context.Background(), signature)

	return nil
}

func GetHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "Cant hash this password"})
	}
	return hashedPassword, nil
}

func CheckRedisToken(id int, token string, name string) error {
	redisClient := redis.GetClient()
	otp, err := redisClient.Get(context.Background(), name+":"+strconv.Itoa(id)).Result()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "Code expired"})
	}

	if otp != token {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "Invalid code"})
	}
	redisClient.Del(context.Background(), name+":"+strconv.Itoa(id))
	return nil
}

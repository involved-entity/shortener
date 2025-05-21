package users

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net/http"
	api "shortener/internal/api"
	conf "shortener/internal/config"
	"shortener/internal/machinery"
	"shortener/internal/redis"
	"strconv"
	"time"

	machineryTasks "github.com/RichardKnop/machinery/v2/tasks"
	"github.com/labstack/echo/v4"
)

func GenerateSecureToken() (string, error) {
	const digits = "0123456789"
	const length = 5
	token := make([]byte, length)

	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			log.Printf("failed to generate token: %v", err)
			return "", errors.New("failed to generate token")
		}
		token[i] = digits[num.Int64()]
	}

	return string(token), nil
}

func CreateAndSendToken(c echo.Context, id uint, email string) error {
	tokenOTP, err := GenerateSecureToken()
	if err != nil {
		return c.JSON(
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

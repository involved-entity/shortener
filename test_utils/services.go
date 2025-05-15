package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	conf "shortener/internal/config"
	"shortener/internal/database"
	r "shortener/internal/redis"
	"strings"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func ExitTest(rClient *redis.Client, exitCode int) {
	rClient.Close()
	conn := database.GetDB()
	tables := []string{
		"users",
		"urls",
		"clicks",
	}

	query := "DROP TABLE IF EXISTS " + strings.Join(tables, ", ") + " CASCADE;"
	err := conn.Exec(query).Error
	if err != nil {
		log.Fatalf("Ошибка при очистке базы данных: %v", err)
	}

	os.Exit(exitCode)
}

func GetJWTForTest(t *testing.T, JWT string) *jwt.Token {
	config := conf.GetConfig()
	parsedToken, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SECRET), nil
	})
	assert.NoError(t, err)
	return parsedToken
}

func GetRedisVarForTestUser(redisName string) string {
	redisClient := r.GetClient()
	val, err := redisClient.Get(context.Background(), redisName+":1").Result()
	if err != nil {
		return ""
	}

	return val
}

func LoginHelper(register, login func(c echo.Context) error) string {
	if err := registerUser(register); err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	token, err := loginUser(login, "testu", "12345678")
	if err != nil {
		log.Fatalf("Failed to login user: %v", err)
	}

	return token
}

func registerUser(register func(c echo.Context) error) error {
	basicTest := BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/register",
		Data:           bytes.NewBuffer([]byte(`{"username": "testu", "email": "test@example.com", "password": "12345678"}`)),
		ExpectedStatus: http.StatusOK,
		Handler:        register,

		NotTestEnv: true,
	}

	basicTest.Execute()

	db := database.GetDB()
	if err := db.Model(&database.User{}).Where("username = ?", "testu").Update("is_verified", true).Error; err != nil {
		return err
	}

	return nil
}

func loginUser(login func(c echo.Context) error, username, password string) (string, error) {
	basicTest := BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/login",
		Data:           bytes.NewBuffer([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))),
		ExpectedStatus: http.StatusOK,
		Handler:        login,

		ServeHTTPMode: true,
		NotTestEnv:    true,
	}

	rec := basicTest.Execute()

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		return "", err
	}

	return response.Data.Token, nil
}

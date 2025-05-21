package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	conf "shortener/internal/config"
	"shortener/internal/database"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func GetJWTForTest(t *testing.T, JWT string) *jwt.Token {
	config := conf.GetConfig()
	parsedToken, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SECRET), nil
	})
	assert.NoError(t, err)
	return parsedToken
}

func LoginHelper(register, login func(c echo.Context) error) string {
	e := echo.New()
	e.POST("/api/register", register)
	e.POST("/api/login", login)

	if err := registerUser(e); err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	token, err := loginUser(e, "testu", "12345678")
	if err != nil {
		log.Fatalf("Failed to login user: %v", err)
	}

	return token
}

func registerUser(e *echo.Echo) error {
	reqBody := `{"username": "testu", "email": "test@example.com", "password": "12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	db := database.GetDB()
	if err := db.Model(&database.User{}).Where("username = ?", "testu").Update("is_verified", true).Error; err != nil {
		return err
	}

	return nil
}

func loginUser(e *echo.Echo, username, password string) (string, error) {
	reqBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

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

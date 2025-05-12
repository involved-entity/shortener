package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"shortener/internal/api"
	conf "shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/redis"
	testutils "shortener/test_utils"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var JWT string

func TestMain(m *testing.M) {
	rClient := testutils.InitTest()

	exitCode := m.Run()

	rClient.Close()
	conn := database.GetDB()
	tables := []string{
		"users",
		"urls",
		"clicks",
	}

	query := "DROP TABLE " + strings.Join(tables, ", ") + ";"
	conn.Exec(query)

	os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(`{"username": "testu", "email": "test@example.com", "password": "12345678"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response api.Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestActivateAccountInvalidCode(t *testing.T) {
	e := echo.New()
	e.POST("/api/verification", ActivateAccount)
	req := httptest.NewRequest(http.MethodPost, "/api/verification", bytes.NewBuffer([]byte(`{
		"id": 1,
		"code": "00000"
	}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestActivateAccount(t *testing.T) {
	e := echo.New()
	redisClient := redis.GetClient()
	config := conf.GetConfig()
	otp, err := redisClient.Get(context.Background(), config.OTP.RedisName+":1").Result()
	if err != nil {
		return
	}
	req := httptest.NewRequest(http.MethodPost, "/api/verification", bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"id": 1,
		"code": "%v"
	}`, otp))))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ActivateAccount(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func TestLogin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer([]byte(`{"username": "testu", "password": "12345678"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Expected 'data' field in response")

	token, ok := data["token"].(string)
	assert.True(t, ok, "Expected 'token' field in data")

	JWT = token
}

func TestResetPassword(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewBuffer([]byte(`{"username": "testu"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ResetPassword(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func TestResetPasswordConfirmInvalidToken(t *testing.T) {
	e := echo.New()
	e.POST("/api/reset-password-confirm", ResetPasswordConfirm)
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password-confirm", bytes.NewBuffer([]byte(`{
		"id": 1,
		"token": "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklz",
		"password": "new-password"
	}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordConfirm(t *testing.T) {
	e := echo.New()
	redisClient := redis.GetClient()
	config := conf.GetConfig()
	token, err := redisClient.Get(context.Background(), config.ResetToken.RedisName+":1").Result()
	if err != nil {
		return
	}
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password-confirm", bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"id": 1,
		"token": "%v",
		"password": "12345678"
	}`, token))))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordConfirm(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMe(t *testing.T) {
	e := echo.New()
	config := conf.GetConfig()
	req := httptest.NewRequest(http.MethodGet, "/api/account", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SECRET), nil
	})
	assert.NoError(t, err)

	c.Set("user", parsedToken)

	err = GetMe(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateAccount(t *testing.T) {
	e := echo.New()
	config := conf.GetConfig()
	req := httptest.NewRequest(http.MethodPatch, "/api/account", bytes.NewBuffer([]byte(`{"email": "newemail@example.com"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SECRET), nil
	})
	assert.NoError(t, err)

	c.Set("user", parsedToken)

	err = UpdateAccount(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"shortener/internal/api"
	conf "shortener/internal/config"
	testUtils "shortener/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var JWT string

func TestMain(m *testing.M) {
	rClient := testUtils.InitTest()

	exitCode := m.Run()

	testUtils.ExitTest(rClient, exitCode)
}

func TestRegister(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/register",
		Data:           bytes.NewBuffer([]byte(`{"username": "testu", "email": "test@example.com", "password": "12345678"}`)),
		ExpectedStatus: http.StatusOK,
		Handler:        Register,
		T:              t,
	}

	rec := basicTest.Execute()

	var response api.Response
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestActivateAccountInvalidCode(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/verification",
		Data:           bytes.NewBuffer([]byte(`{"id": 1, "code": "00000"}`)),
		ExpectedStatus: http.StatusBadRequest,
		Handler:        ActivateAccount,
		T:              t,

		ServeHTTPMode: true,
	}

	basicTest.Execute()
}

func TestActivateAccount(t *testing.T) {
	otp := testUtils.GetRedisVarForTestUser(conf.GetConfig().OTP.RedisName)

	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/verification",
		Data:           bytes.NewBuffer([]byte(fmt.Sprintf(`{"id": 1, "code": "%v"}`, otp))),
		ExpectedStatus: http.StatusAccepted,
		Handler:        ActivateAccount,
		T:              t,
	}

	basicTest.Execute()
}

func TestLogin(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/login",
		Data:           bytes.NewBuffer([]byte(`{"username": "testu", "password": "12345678"}`)),
		ExpectedStatus: http.StatusOK,
		Handler:        Login,
		T:              t,
	}

	rec := basicTest.Execute()

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Expected 'data' field in response")

	token, ok := data["token"].(string)
	assert.True(t, ok, "Expected 'token' field in data")

	JWT = token
}

func TestResetPassword(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/reset-password",
		Data:           bytes.NewBuffer([]byte(`{"username": "testu"}`)),
		ExpectedStatus: http.StatusAccepted,
		Handler:        ResetPassword,
		T:              t,
	}

	basicTest.Execute()
}

func TestResetPasswordConfirmInvalidToken(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method: http.MethodPost,
		Url:    "/api/reset-password-confirm",
		Data: bytes.NewBuffer([]byte(`{
			"id": 1,
			"token": "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklz",
			"password": "new-password"
		}`)),
		ExpectedStatus: http.StatusBadRequest,
		Handler:        ResetPasswordConfirm,
		T:              t,

		ServeHTTPMode: true,
	}

	basicTest.Execute()
}

func TestResetPasswordConfirm(t *testing.T) {
	token := testUtils.GetRedisVarForTestUser(conf.GetConfig().ResetToken.RedisName)

	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/reset-password-confirm",
		Data:           bytes.NewBuffer([]byte(fmt.Sprintf(`{"id": 1, "token": "%v", "password": "12345678"}`, token))),
		ExpectedStatus: http.StatusOK,
		Handler:        ResetPasswordConfirm,
		T:              t,
	}

	basicTest.Execute()
}

func TestGetMe(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodGet,
		Url:            "/api/account",
		Data:           &bytes.Buffer{},
		ExpectedStatus: http.StatusOK,
		Handler:        GetMe,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

func TestUpdateAccount(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPatch,
		Url:            "/api/account",
		Data:           bytes.NewBuffer([]byte(`{"email": "newemail@example.com"}`)),
		ExpectedStatus: http.StatusOK,
		Handler:        UpdateAccount,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

package urls

import (
	"bytes"
	"net/http"
	"shortener/internal/api/users"
	testUtils "shortener/test_utils"
	"testing"
)

var JWT string

func TestMain(m *testing.M) {
	rClient := testUtils.InitTest()

	JWT = testUtils.LoginHelper(users.Register, users.Login)

	exitCode := m.Run()

	testUtils.ExitTest(rClient, exitCode)
}

func TestSaveURL(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodPost,
		Url:            "/api/urls",
		Data:           bytes.NewBuffer([]byte(`{"originalURL":"https://example.com","shortCode":"sh"}`)),
		ExpectedStatus: http.StatusOK,
		Handler:        SaveURL,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

func TestGetURL(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodGet,
		Url:            "/_/sh",
		Data:           &bytes.Buffer{},
		ExpectedStatus: http.StatusPermanentRedirect,
		Handler:        GetURL,
		T:              t,

		ServeHTTPMode: true,
		HandlerUrl:    "/_/:shortCode",
	}

	basicTest.Execute()
}

func TestDeleteURL(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodDelete,
		Url:            "/api/urls/sh",
		Data:           &bytes.Buffer{},
		ExpectedStatus: http.StatusNoContent,
		Handler:        DeleteURL,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

func TestGetMyURLs(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodGet,
		Url:            "/api/urls",
		Data:           &bytes.Buffer{},
		ExpectedStatus: http.StatusOK,
		Handler:        GetMyURLs,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

func TestGetURLClicks(t *testing.T) {
	basicTest := testUtils.BasicTest{
		Method:         http.MethodGet,
		Url:            "/api/clicks/sh",
		Data:           &bytes.Buffer{},
		ExpectedStatus: http.StatusOK,
		Handler:        GetURLClicks,
		T:              t,

		JWT: JWT,
	}

	basicTest.Execute()
}

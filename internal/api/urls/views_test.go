package urls

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"shortener/internal/api/users"
	"shortener/internal/database"
	testutils "shortener/test_utils"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var JWT string

func TestMain(m *testing.M) {
	rClient := testutils.InitTest()

	JWT = testutils.LoginHelper(users.Register, users.Login)

	exitCode := m.Run()

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

func TestSaveURL(t *testing.T) {
	e := echo.New()
	urlDTO := `{"originalURL":"https://example.com","shortCode":"sh"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(urlDTO))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken := testutils.GetJWTForTest(t, JWT)
	c.Set("user", parsedToken)

	err := SaveURL(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetURL(t *testing.T) {
	e := echo.New()
	e.GET("/_/:shortCode", GetURL)
	req := httptest.NewRequest(http.MethodGet, "/_/sh", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken := testutils.GetJWTForTest(t, JWT)
	c.Set("user", parsedToken)

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
}

func TestDeleteURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/urls/sh", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken := testutils.GetJWTForTest(t, JWT)
	c.Set("user", parsedToken)

	err := DeleteURL(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestGetMyURLs(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken := testutils.GetJWTForTest(t, JWT)
	c.Set("user", parsedToken)

	err := GetMyURLs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetURLClicks(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/clicks/sh", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	parsedToken := testutils.GetJWTForTest(t, JWT)
	c.Set("user", parsedToken)

	err := GetURLClicks(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

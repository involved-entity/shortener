package testutils

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type BasicTest struct {
	Method         string
	Url            string
	Data           *bytes.Buffer
	ExpectedStatus int
	Handler        func(e echo.Context) error
	T              *testing.T

	JWT           string
	ServeHTTPMode bool
	HandlerUrl    string
	NotTestEnv    bool
}

func (t BasicTest) Execute() *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(t.Method, t.Url, t.Data)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	if t.ServeHTTPMode {
		val := reflect.ValueOf(e)
		method := val.MethodByName(t.Method)
		url := ternary(t.HandlerUrl != "", t.HandlerUrl, t.Url).(string)
		argVals := []reflect.Value{reflect.ValueOf(url), reflect.ValueOf(t.Handler)}
		method.Call(argVals)

		e.ServeHTTP(rec, req)
	} else {
		c := e.NewContext(req, rec)

		if t.JWT != "" {
			parsedToken := GetJWTForTest(t.T, t.JWT)
			c.Set("user", parsedToken)
		}

		err := t.Handler(c)
		if !t.NotTestEnv {
			assert.NoError(t.T, err)
		}
	}

	if !t.NotTestEnv {
		assert.Equal(t.T, t.ExpectedStatus, rec.Code)
	}

	return rec
}

func ternary(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

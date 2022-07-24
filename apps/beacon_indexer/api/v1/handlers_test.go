package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	E *echo.Echo
}

func (s *HandlersTestSuite) SetupTest() {
	s.E = echo.New()
	s.E = Routes(s.E)
}

func (s *HandlersTestSuite) TestSetLevel() {
	s.getRequest("http://localhost/log/0", http.StatusOK)
}

func (s *HandlersTestSuite) getRequest(reqPath string, httpCode int) {
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	rec := httptest.NewRecorder()
	s.E.ServeHTTP(rec, req)
	s.Equal(httpCode, rec.Code)

	return
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

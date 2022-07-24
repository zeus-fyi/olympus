package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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

func (s *HandlersTestSuite) TestAdminCfg() {
	ll := zerolog.DebugLevel
	nvSize, nbSize := 100, 500
	timeout := time.Second * 30
	adminCfg := AdminConfig{
		LogLevel:                   &ll,
		ValidatorBatchSize:         &nvSize,
		ValidatorBalancesBatchSize: &nbSize,
		ValidatorBalancesTimeout:   &timeout,
	}

	adminReq := AdminConfigRequest{adminCfg}
	resp := s.postAdminRequest(adminReq, http.StatusOK)
	s.Assert().NotEmpty(resp)
}

func (s *HandlersTestSuite) TestSetLevel() {
	s.getRequest("http://localhost/health", http.StatusOK)
}

func (s *HandlersTestSuite) getRequest(reqPath string, httpCode int) {
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	rec := httptest.NewRecorder()
	s.E.ServeHTTP(rec, req)
	s.Equal(httpCode, rec.Code)
	return
}

type TestResponse struct {
	logs []byte
}

func (s *HandlersTestSuite) postAdminRequest(postRequest AdminConfigRequest, httpCode int) TestResponse {
	podActionRequestPayload, err := json.Marshal(postRequest)
	s.Assert().Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/admin", strings.NewReader(string(podActionRequestPayload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	s.E.ServeHTTP(rec, req)
	s.Equal(httpCode, rec.Code)

	var tr TestResponse
	tr.logs = rec.Body.Bytes()

	return tr
}
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

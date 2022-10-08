package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (s *HandlersTestSuite) TestAdminRedisCfg() {

	adminReq := AdminRedisConfigRequest{
		Addr:   "localhost:6379",
		OsEnv:  "REDIS",
		UseEnv: false,
	}

	resp := s.postAdminRequest(adminReq, "debug/redis", http.StatusOK)
	s.Assert().NotEmpty(resp)
}

//func (s *HandlersTestSuite) TestAdminCfg() {
//	ll := zerolog.DebugLevel
//	nvSize, nbSize := 100, 500
//	timeout := time.Second * 30
//	adminCfg := AdminConfig{
//		LogLevel:                   &ll,
//		ValidatorBatchSize:         &nvSize,
//		ValidatorBalancesBatchSize: &nbSize,
//		ValidatorBalancesTimeout:   &timeout,
//	}
//
//	adminReq := AdminConfigRequest{adminCfg}
//	resp := s.postAdminRequest(adminReq, "debug/db/config", http.StatusOK)
//	s.Assert().NotEmpty(resp)
//}
//
//func (s *HandlersTestSuite) TestAdminDBCfg() {
//
//	maxConn := int32(10)
//	adminReq := AdminDBConfigRequest{
//		postgres.ConfigChangePG{
//			MaxConns: &maxConn,
//		},
//	}
//	resp := s.postAdminRequest(adminReq, "admin", http.StatusOK)
//	s.Assert().NotEmpty(resp)
//}

func (s *HandlersTestSuite) TestGetAdminInfo() {
	s.getRequest("http://localhost/admin", http.StatusOK)
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

func (s *HandlersTestSuite) postAdminRequest(postRequest AdminRedisConfigRequest, endpoint string, httpCode int) TestResponse {
	podActionRequestPayload, err := json.Marshal(postRequest)

	s.Assert().Nil(err)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/"+endpoint, strings.NewReader(string(podActionRequestPayload)))
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

package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/env"
)

type ValidatorBalancesTestSuite struct {
	env.StagingPrototypeTest
	E *echo.Echo
}

func (s *ValidatorBalancesTestSuite) TestValidatorBalancesRequest() {
	s.E = echo.New()
	s.E = Routes(s.E)

	le, he := 0, 3
	vbr := ValidatorBalancesRequest{
		LowerEpoch:       le,
		HigherEpoch:      he,
		ValidatorIndexes: []int64{1, 2, 3, 4, 5},
	}

	tr := s.postValidatorBalancesRequest(vbr, "v1/validator_balances", 200)
	s.Assert().NotEmpty(tr.logs)
}

func (s *ValidatorBalancesTestSuite) TestValidatorBalancesSumRequest() {
	s.E = echo.New()
	s.E = Routes(s.E)

	le, he := 0, 3
	vbr := ValidatorBalancesRequest{
		LowerEpoch:       le,
		HigherEpoch:      he,
		ValidatorIndexes: []int64{1, 2, 3, 4, 5},
	}

	tr := s.postValidatorBalancesRequest(vbr, "v1/validator_balances_sums", 200)
	s.Assert().NotEmpty(tr.logs)
}

func (s *ValidatorBalancesTestSuite) getRequest(reqPath string, httpCode int) {
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	rec := httptest.NewRecorder()
	s.E.ServeHTTP(rec, req)
	s.Equal(httpCode, rec.Code)
	return
}

func (s *ValidatorBalancesTestSuite) postValidatorBalancesRequest(postRequest ValidatorBalancesRequest, endpoint string, httpCode int) TestResponse {
	podActionRequestPayload, err := json.Marshal(postRequest)

	s.Assert().Nil(err)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/"+endpoint, strings.NewReader(string(podActionRequestPayload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	req.Header.Set(echo.HeaderAuthorization, bearer)

	rec := httptest.NewRecorder()
	s.E.ServeHTTP(rec, req)
	s.Equal(httpCode, rec.Code)

	var tr TestResponse
	tr.logs = rec.Body.Bytes()

	return tr
}
func TestValidatorBalancesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorBalancesTestSuite))
}

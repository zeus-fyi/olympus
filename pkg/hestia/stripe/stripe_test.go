package hestia_stripe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type StripeTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *StripeTestSuite) TestSendTestEmail() {
	s.InitLocalConfigs()
}

func TestStripeTestSuite(t *testing.T) {
	suite.Run(t, new(StripeTestSuite))
}

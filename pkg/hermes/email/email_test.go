package hermes_email_notifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aws_aegis_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type EmailTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *EmailTestSuite) SendTestEmail() {
	auth := aws_aegis_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySES,
		SecretKey: s.Tc.AwsSecretKeySES,
	}
	h := InitHermesEmailNotifications(ctx, auth)
	s.Require().Nil(h.Client)
	r, err := h.SendEmailTo(ctx, []string{"alex@zeus.fyi"})
	s.Require().Nil(err)
	s.Require().NotNil(r)
}

func TestEmailTestSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

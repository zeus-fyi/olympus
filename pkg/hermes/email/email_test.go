package hermes_email_notifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aws_aegis_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type EmailTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *EmailTestSuite) TestSendTestEmail() {
	s.InitLocalConfigs()
	auth := aws_aegis_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySES,
		SecretKey: s.Tc.AwsSecretKeySES,
	}
	h := InitHermesSESEmailNotifications(ctx, auth)
	s.Require().NotNil(h.SES)
	us := create_org_users.UserSignup{
		FirstName:        "",
		LastName:         "",
		EmailAddress:     "ageorge010@vt.edu",
		Password:         "",
		VerifyEmailToken: "abc123",
	}
	r, err := h.SendSESEmailVerifyRequest(ctx, us)
	s.Require().Nil(err)
	s.Require().NotNil(r)
}

func TestEmailTestSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

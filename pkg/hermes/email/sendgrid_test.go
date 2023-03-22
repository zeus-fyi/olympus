package hermes_email_notifications

import create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"

func (s *EmailTestSuite) TestSendGridEmail() {
	InitHermesSendGridClient(ctx, s.Tc.SendGridAPIKey)
	s.Require().NotNil(Hermes.SendGrid)
	us := create_org_users.UserSignup{
		FirstName:        "alex",
		LastName:         "g",
		EmailAddress:     "ageorge010@vt.edu",
		Password:         "",
		VerifyEmailToken: "abc123",
	}
	r, err := Hermes.SendSendGridEmailVerifyRequest(ctx, us)
	s.Require().Nil(err)
	s.Require().NotNil(r)
}

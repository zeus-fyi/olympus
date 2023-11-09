package hermes_email_notifications

func (s *EmailTestSuite) TestNewGmail() {

	em := "alex@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)

	em = "support@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)
}

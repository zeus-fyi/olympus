package hermes_email_notifications

func (s *EmailTestSuite) TestNewGmail() {

	em := "alex@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)

	em = "support@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)

	NewGmailServiceClient(ctx, s.Tc.GcpAuthJson, em)
}

func (s *EmailTestSuite) TestNewGmailWorker() {
	em := "alex@zeus.fyi"
	gs := NewGmailServiceClient(ctx, s.Tc.GcpAuthJson, em)

	//gs.ReadEmails(em)
	//em = "support@zeus.fyi"
	gs.ReadEmails(em)
}

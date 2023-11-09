package hermes_email_notifications

func (s *EmailTestSuite) TestNewGmail() {

	gm, err := NewGmail(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(gm)

}

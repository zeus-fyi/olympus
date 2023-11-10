package hermes_email_notifications

import (
	"fmt"
)

func (s *EmailTestSuite) TestNewGmail() {

	em := "alex@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)

	em = "support@zeus.fyi"
	NewGmail(ctx, s.Tc.GcpAuthJson, em)

	NewGmailServiceClient(ctx, s.Tc.GcpAuthJson, em)
}

func (s *EmailTestSuite) TestNewGmailWorker() {
	em := "ai@zeus.fyi"
	gs := NewGmailServiceClient(ctx, s.Tc.GcpAuthJson, em)

	emailContents, err := gs.GetReadEmails(em, 1)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for _, emailContent := range emailContents {
		fmt.Println("=========================================")
		fmt.Println("MsgId: ", emailContent.MsgId)
		fmt.Println("Email: ", emailContent.From)
		fmt.Println("Subject: ", emailContent.Subject)

		fmt.Println("Body: ", stripHTML(emailContent.Body))
	}
}

package hermes_email_notifications

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	email_templates "github.com/zeus-fyi/olympus/pkg/hermes/email/templates"
)

func InitHermesSendGridClient(ctx context.Context, apiKey string) {
	client := sendgrid.NewSendClient(apiKey)
	Hermes.SendGrid = client
	return
}

func (h *HermesEmailNotifications) SendSendGridEmailVerifyRequest(ctx context.Context, us create_org_users.UserSignup) (*rest.Response, error) {
	from := mail.NewEmail("Zeus Cloud", "alex@zeus.fyi")
	subject := "Verify Your Email at Zeus Cloud"
	to := mail.NewEmail(fmt.Sprintf("%s", us.FirstName), us.EmailAddress)
	htmlContent := email_templates.VerifyEmailHTML(us.VerifyEmailToken)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	resp, err := h.SendGrid.Send(message)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("HermesEmailNotifications: SendSendGridEmailVerifyRequest: error")
		return nil, err
	}
	return resp, err
}

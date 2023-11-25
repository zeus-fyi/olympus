package hermes_email_notifications

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	email_templates "github.com/zeus-fyi/olympus/pkg/hermes/email/templates"
	aws_aegis_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var Hermes = HermesEmailNotifications{}

type HermesEmailNotifications struct {
	SendGrid *sendgrid.Client
	SES      *sesv2.Client
}

func InitHermesSESEmailNotifications(ctx context.Context, a aws_aegis_auth.AuthAWS) HermesEmailNotifications {
	cfg, err := a.CreateConfig(ctx)
	if err != nil {
		log.Err(err).Msg("InitSES: error loading config")
		panic(err)
	}
	client := sesv2.NewFromConfig(cfg)
	return HermesEmailNotifications{SES: client}
}

func (h *HermesEmailNotifications) SendSESEmailVerifyRequest(ctx context.Context, us create_org_users.UserSignup) (*sesv2.SendEmailOutput, error) {
	html := email_templates.VerifyEmailHTML(us.VerifyEmailToken)
	params := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(html),
						Charset: aws.String("utf-8"),
					},
				},
				Subject: &types.Content{
					Data:    aws.String("Verify Your Email at Zeus Cloud"),
					Charset: aws.String("utf-8"),
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{us.EmailAddress},
		},
		FromEmailAddress: aws.String("alex@zeus.fyi"),
	}
	resp, err := h.SES.SendEmail(ctx, params)
	if err != nil {
		log.Ctx(ctx).Info().Interface("resp", resp).Err(err).Msg("HermesEmailNotifications: SendEmailVerify: error")
		return nil, err
	}
	return resp, err
}

package hermes_email_notifications

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/rs/zerolog/log"
	aws_aegis_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var Hermes = HermesEmailNotifications{}

type HermesEmailNotifications struct {
	*sesv2.Client
}

func InitHermesEmailNotifications(ctx context.Context, a aws_aegis_auth.AuthAWS) HermesEmailNotifications {
	cfg, err := a.CreateConfig(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InitSES: error loading config")
		panic(err)
	}
	client := sesv2.NewFromConfig(cfg)
	return HermesEmailNotifications{Client: client}
}

func (h *HermesEmailNotifications) SendEmailTo(ctx context.Context, toEmails []string) (*sesv2.SendEmailOutput, error) {

	// TODO token creation
	params := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: nil,
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data:    nil,
						Charset: nil,
					},
					Text: &types.Content{
						Data:    nil,
						Charset: nil,
					},
				},
				Subject: &types.Content{
					Data:    aws.String("Email Verification"),
					Charset: nil,
				},
			},
			Template: nil,
		},
		ConfigurationSetName: nil,
		Destination: &types.Destination{
			BccAddresses: nil,
			CcAddresses:  nil,
			ToAddresses:  toEmails,
		},
		FeedbackForwardingEmailAddress:            nil,
		FeedbackForwardingEmailAddressIdentityArn: nil,
		FromEmailAddress:                          nil,
		FromEmailAddressIdentityArn:               nil,
		ListManagementOptions:                     nil,
		ReplyToAddresses:                          nil,
	}
	resp, err := h.SendEmail(ctx, params)
	if err != nil {
		log.Ctx(ctx).Info().Interface("resp", resp).Err(err).Msg("HermesEmailNotifications: SendEmailVerify: error")
		panic(err)
	}
	return resp, err
}

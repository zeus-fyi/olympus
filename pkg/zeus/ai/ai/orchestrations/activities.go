package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

type ZeusAiPlatformActivities struct {
	kronos_helix.ActivityDefinition
}

func NewZeusAiPlatformActivities() ZeusAiPlatformActivities {
	return ZeusAiPlatformActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *ZeusAiPlatformActivities) GetActivities() ActivitiesSlice {
	ka := kronos_helix.NewKronosActivities()
	actSlice := []interface{}{h.AiTask, h.SaveAiTaskResponse, h.SendTaskResponseEmail, h.InsertEmailIfNew,
		h.InsertAiResponse, h.InsertTelegramMessageIfNew,
	}
	return append(actSlice, ka.GetActivities()...)
}

func (h *ZeusAiPlatformActivities) AiTask(ctx context.Context, ou org_users.OrgUser, msg hermes_email_notifications.EmailContents) (openai.ChatCompletionResponse, error) {
	//task := "write a bullet point summary of the email contents and suggest some responses if applicable. write your reply as html formatted\n"
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful bot that reads email contents and provides a bullet point summary and then suggest well thought out responses and that aren't overly formal or stiff in tone and you always write your reply as well formatted html that is easy to read.",
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
			Messages: []openai.ChatCompletionMessage{
				systemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: hermes_email_notifications.GenerateAiRequest(msg.Body, msg),
					Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
				},
			},
		},
	)
	return resp, err
}

func (h *ZeusAiPlatformActivities) SaveAiTaskResponse(ctx context.Context, ou org_users.OrgUser, resp openai.ChatCompletionResponse) error {
	err := hera_openai.HeraOpenAI.RecordUIChatRequestUsage(ctx, ou, resp)
	if err != nil {
		log.Err(err).Msg("SaveAiTaskResponse: RecordUIChatRequestUsage failed")
		return nil
	}
	return nil
}

func (h *ZeusAiPlatformActivities) SendTaskResponseEmail(ctx context.Context, email string, resp openai.ChatCompletionResponse) error {
	content := ""
	for _, msg := range resp.Choices {
		// Remove markdown code block characters
		line := strings.Replace(msg.Message.Content, "```", "", -1)

		//// Escape any HTML special characters to prevent XSS or other issues
		//line = html.EscapeString(line)

		// Add the line break for proper formatting in HTML
		content += line
	}

	if len(content) == 0 {
		return nil
	}
	_, err := hermes_email_notifications.Hermes.SendAITaskResponse(ctx, email, content)
	if err != nil {
		log.Err(err).Msg("SendTaskResponseEmail: SendAITaskResponse failed")
		return err
	}
	return nil
}

func (h *ZeusAiPlatformActivities) InsertEmailIfNew(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

func (h *ZeusAiPlatformActivities) InsertTelegramMessageIfNew(ctx context.Context, msg any) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewTgMessages(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

func (h *ZeusAiPlatformActivities) InsertAiResponse(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

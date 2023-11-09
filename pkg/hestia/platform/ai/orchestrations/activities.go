package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

type HestiaAiPlatformActivities struct {
	kronos_helix.KronosActivities
}

func NewHestiaAiPlatformActivities() HestiaAiPlatformActivities {
	return HestiaAiPlatformActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaAiPlatformActivities) GetActivities() ActivitiesSlice {
	actSlice := []interface{}{h.AiTask, h.SaveAiTaskResponse, h.SendTaskResponseEmail}
	return actSlice
}

func (h *HestiaAiPlatformActivities) AiTask(ctx context.Context, ou org_users.OrgUser, content string) (openai.ChatCompletionResponse, error) {
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
					Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
				},
			},
		},
	)
	return resp, err
}

func (h *HestiaAiPlatformActivities) SaveAiTaskResponse(ctx context.Context, ou org_users.OrgUser, resp openai.ChatCompletionResponse) error {
	err := hera_openai.HeraOpenAI.RecordUIChatRequestUsage(ctx, ou, resp)
	if err != nil {
		log.Err(err).Msg("SaveAiTaskResponse: RecordUIChatRequestUsage failed")
		return err
	}
	return nil
}

func (h *HestiaAiPlatformActivities) SendTaskResponseEmail(ctx context.Context, email string, resp openai.ChatCompletionResponse) error {
	content := ""
	for _, msg := range resp.Choices {
		content += msg.Message.Content + "\n"
	}
	_, err := hermes_email_notifications.Hermes.SendAITaskResponse(ctx, email, content)
	if err != nil {
		log.Err(err).Msg("SendTaskResponseEmail: SendAITaskResponse failed")
		return err
	}
	return nil
}

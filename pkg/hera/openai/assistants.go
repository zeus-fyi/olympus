package hera_openai

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

func CreateOrUpdateAssistant(ctx context.Context, oac OpenAI, assistant *openai.Assistant) (*openai.Assistant, error) {
	if assistant == nil {
		return nil, nil
	}
	if assistant.ID != "" {
		return UpdateAssistant(ctx, oac, assistant)
	}
	resp, err := oac.CreateAssistant(ctx, openai.AssistantRequest{
		Model:        assistant.Model,
		Name:         assistant.Name,
		Description:  assistant.Description,
		Instructions: assistant.Instructions,
		Tools:        []openai.AssistantTool{},
	})
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("failed to create assistant")
		return nil, err
	}
	return &resp, nil
}

func UpdateAssistant(ctx context.Context, oac OpenAI, assistant *openai.Assistant) (*openai.Assistant, error) {
	if assistant == nil {
		return nil, nil
	}
	resp, err := oac.ModifyAssistant(ctx, assistant.ID, openai.AssistantRequest{
		Model:        assistant.Model,
		Name:         assistant.Name,
		Description:  assistant.Description,
		Instructions: assistant.Instructions,
	})
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("failed to create assistant")
		return nil, err
	}
	return &resp, nil
}

func ListAssistants(ctx context.Context, oac OpenAI) ([]openai.Assistant, error) {
	resp, err := oac.ListAssistants(ctx, nil, nil, nil, nil)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("failed to create assistant")
		return nil, err
	}
	return resp.Assistants, nil
}

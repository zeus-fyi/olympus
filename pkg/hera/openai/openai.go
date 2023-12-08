package hera_openai

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

var HeraOpenAI OpenAI

var maxTokensByModel = map[string]int{
	openai.GPT4TurboPreview:   40000,
	openai.GPT3TextDavinci003: 2048,
	openai.GPT3TextDavinci002: 2048,
	openai.GPT3TextDavinci001: 2048,
	openai.GPT3TextAda001:     2048,
	openai.GPT3TextBabbage001: 2048,
}

type OpenAI struct {
	*openai.Client
}

func InitHeraOpenAI(bearer string) {
	HeraOpenAI = OpenAI{}
	HeraOpenAI.Client = openai.NewClient(bearer)
}

type OpenAIParams struct {
	Model     string
	MaxTokens int
	Prompt    string
}

func (ai *OpenAI) RecordUIChatRequestUsage(ctx context.Context, ou org_users.OrgUser, params openai.ChatCompletionResponse, prompt []byte) error {
	log.Info().Interface("params", params).Msg("RecordUIChatRequestUsage")
	_, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, ou, params, prompt)
	if err != nil {
		log.Err(err).Interface("params", params).Msg("RecordUIChatRequestUsage")
		return err
	}
	return nil
}

func (ai *OpenAI) MakeCodeGenRequest(ctx context.Context, ou org_users.OrgUser, params OpenAIParams) (openai.CompletionResponse, error) {
	params.Model = openai.GPT4TurboPreview
	params.MaxTokens = maxTokensByModel[params.Model]
	req := openai.CompletionRequest{
		Model:     params.Model,
		MaxTokens: params.MaxTokens,
		Prompt:    params.Prompt,
		User:      fmt.Sprintf("%d", ou.UserID),
	}

	resp, err := ai.CreateCompletion(ctx, req)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("MakeCodeGenRequest")
		return resp, err
	}
	err = hera_openai_dbmodels.InsertCompletionResponse(ctx, ou, resp)
	if err != nil {
		log.Err(err)
		return resp, err
	}
	return resp, err
}

func (ai *OpenAI) MakeCodeGenRequestV2(ctx context.Context, ou org_users.OrgUser, params OpenAIParams) (openai.ChatCompletionResponse, error) {
	systemMessage := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: "You are a helpful bot that analyzes the context, filepaths, and content of supplied code references and generates code from example functions, code references, and other guidance." +
			" You respond only with code, and you are not a chatbot",
		Name: fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	resp, err := ai.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
			Messages: []openai.ChatCompletionMessage{
				systemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: params.Prompt,
					Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
				},
			},
		},
	)
	return resp, err
}

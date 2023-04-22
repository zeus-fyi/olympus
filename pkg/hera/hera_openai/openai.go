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

func (ai *OpenAI) RecordUIChatRequestUsage(ctx context.Context, ou org_users.OrgUser, params openai.ChatCompletionResponse) error {
	cmpChoices := make([]openai.CompletionChoice, len(params.Choices))
	for i, ch := range params.Choices {
		cmpChoices[i].Text = ch.Message.Content
		cmpChoices[i].Index = ch.Index
		cmpChoices[i].FinishReason = ch.FinishReason
	}

	resp := openai.CompletionResponse{
		ID:      params.ID,
		Object:  params.Object,
		Created: params.Created,
		Model:   params.Model,
		Usage:   params.Usage,
		Choices: cmpChoices,
	}
	err := hera_openai_dbmodels.InsertCompletionResponse(ctx, ou, resp)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return err
}

func (ai *OpenAI) MakeCodeGenRequest(ctx context.Context, ou org_users.OrgUser, params OpenAIParams) (openai.CompletionResponse, error) {
	if len(params.Model) <= 0 {
		params.Model = openai.GPT3TextDavinci003
	}

	req := openai.CompletionRequest{
		Model:     params.Model,
		MaxTokens: params.MaxTokens,
		Prompt:    params.Prompt,
		User:      fmt.Sprintf("%d", ou.UserID),
	}

	resp, err := ai.CreateCompletion(ctx, req)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("orgUser", ou).Msg("MakeCodeGenRequest")
		return resp, err
	}
	err = hera_openai_dbmodels.InsertCompletionResponse(ctx, ou, resp)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return resp, err
	}
	return resp, err
}

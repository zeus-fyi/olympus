package openai

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	gogpt "github.com/sashabaranov/go-gpt3"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

var HeraOpenAI OpenAI

var maxTokensByModel = map[string]int{
	gogpt.GPT3TextDavinci003: 2048,
	gogpt.GPT3TextDavinci002: 2048,
	gogpt.GPT3TextDavinci001: 2048,
	gogpt.GPT3TextAda001:     2048,
	gogpt.GPT3TextBabbage001: 2048,
}

type OpenAI struct {
	*gogpt.Client
}

func InitHeraOpenAI(bearer string) {
	HeraOpenAI = OpenAI{}
	HeraOpenAI.Client = gogpt.NewClient(bearer)
}

type OpenAIParams struct {
	Model     string
	MaxTokens int
	Prompt    string
}

func (ai *OpenAI) MakeCodeGenRequest(ctx context.Context, ou org_users.OrgUser, params OpenAIParams) (gogpt.CompletionResponse, error) {
	if len(params.Model) <= 0 {
		params.Model = gogpt.GPT3TextDavinci003
	}

	req := gogpt.CompletionRequest{
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

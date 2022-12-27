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

const (
	codexMaxTokens = 8000
)

type OpenAI struct {
	*gogpt.Client
}

func InitHeraOpenAI(bearer string) {
	HeraOpenAI = OpenAI{}
	HeraOpenAI.Client = gogpt.NewClient(bearer)
}

func (ai *OpenAI) MakeCodeGenRequest(ctx context.Context, prompt string, ou org_users.OrgUser) (gogpt.CompletionResponse, error) {
	req := gogpt.CompletionRequest{
		Model:     gogpt.CodexCodeDavinci002,
		MaxTokens: codexMaxTokens,
		Prompt:    prompt,
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
	}
	return resp, err
}

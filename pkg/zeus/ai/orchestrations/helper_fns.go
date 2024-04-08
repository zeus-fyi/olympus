package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

func GetPandoraMessages(ctx context.Context, ou org_users.OrgUser, groupPrefix string) ([]hera_search.TelegramMessage, error) {
	token, err := aws_secrets.GetTelegramToken(ctx, ou.OrgID)
	if err != nil {
		aws_secrets.SecretsCache.Delete("telegram_token")
		log.Err(err).Msg("Zeus: GetTelegramToken")
		return nil, err
	}

	var msgs []hera_search.TelegramMessage
	apiReq := &iris_api_requests.ApiProxyRequest{
		Url:             "https://pandora.zeus.fyi",
		PayloadTypeREST: "POST",
		Payload: echo.Map{
			"group": groupPrefix,
			"token": token,
		},
		IsInternal: true,
	}

	res := resty_base.GetBaseRestyClient(apiReq.Url, artemis_orchestration_auth.Bearer)
	resp, err := res.R().SetBody(&apiReq.Payload).SetResult(&msgs).Post("msgs")
	if err != nil {
		log.Err(err).Msg("Zeus: GetPandoraMessages")
		return nil, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		if err != nil {
			err = fmt.Errorf("Zeus: GetPandoraMessages: failed to relay api request: status code %d", resp.StatusCode())
		}
		return nil, err
	}
	return msgs, nil
}

type TokenCountsEstimate struct {
	Count int `json:"count"`
}

func GetTokenCountEstimate(ctx context.Context, model, text string) (int, error) {
	if len(model) == 0 {
		model = "gpt-4"
	}
	if strings.HasPrefix(model, "gpt-4") {
		model = "gpt-4"
	}
	if strings.HasPrefix(model, "gpt-3.5") {
		model = "gpt-3.5-turbo"
	}
	var tc TokenCountsEstimate
	apiReq := &iris_api_requests.ApiProxyRequest{
		Url:             "https://pandora.zeus.fyi",
		PayloadTypeREST: "POST",
		Payload: echo.Map{
			"model": model,
			"text":  text,
		},
		IsInternal: true,
	}
	res := resty_base.GetBaseRestyClient(apiReq.Url, api_auth_temporal.Bearer)
	resp, err := res.R().SetBody(&apiReq.Payload).SetResult(&tc).Post("tokenize")
	if err != nil {
		log.Err(err).Interface("&apiReq.Payload)", &apiReq.Payload).Msg("Zeus: GetTokenCountEstimate")
		return -1, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		if err == nil {
			err = fmt.Errorf("GetTokenCountEstimate: failed to relay api request: status code %d", resp.StatusCode())
		}
		log.Err(err).Interface("&apiReq.Payload)", &apiReq.Payload).Msg("Zeus: GetTokenCountEstimate")
		return -1, err
	}
	return tc.Count, nil
}

func AiTelegramTask(ctx context.Context, ou org_users.OrgUser, msgs []hera_search.SearchResult, params hera_search.AiSearchParams) (openai.ChatCompletionResponse, error) {
	systemMessage := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: "You are a helpful sales bot that reads telegram messages and analyzes them for products, services," +
			" or sales related questions and then summarizes and groups the users and predicts which users which need " +
			"to be persuaded more and ranks them by importance. You then generate customer profiles with best effort, " +
			"and predicts which pain points are more relevant, and which talking points aren't interesting or relevant, " +
			"or are confusing and then suggest well thought out next steps the sales team should take. " +
			"You aren't overly formal or stiff in tone",
		Name: fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	if params.Retrieval.RetrievalPrompt != nil && len(*params.Retrieval.RetrievalPrompt) > 0 {
		systemMessage.Content = *params.Retrieval.RetrievalPrompt
	}
	output := hera_search.FormatSearchMessagesForAi(msgs)

	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: Gpt4JsonModel,
			Messages: []openai.ChatCompletionMessage{
				systemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: output,
					Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
				},
			},
		},
	)
	return resp, err
}

func AiAggregateTask(ctx context.Context, ou org_users.OrgUser, msgs []hera_search.SearchResult, params hera_search.AiSearchParams) (openai.ChatCompletionResponse, error) {
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: aws.StringValue(params.Retrieval.RetrievalPrompt),
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	output := hera_search.FormatSearchMessagesForAi(msgs)
	tc, err := GetTokenCountEstimate(ctx, "gpt-4", output)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	if tc > 8000 {
		return openai.ChatCompletionResponse{}, fmt.Errorf("token count too high: %d", tc)
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4",
			Messages: []openai.ChatCompletionMessage{
				systemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: output,
					Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
				},
			},
		},
	)
	return resp, err
}

package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_access_keygen "github.com/zeus-fyi/olympus/hestia/web/access"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

var cah = cache.New(time.Hour, cache.DefaultExpiration)

func GetTelegramToken(ctx context.Context) (string, error) {
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, hestia_access_keygen.FormatSecret(internalOrgID))
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return "", err
	}
	m := make(map[string]hestia_access_keygen.SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return "", err
	}

	token := ""
	tv, ok := cah.Get("telegram_token")
	if ok {
		token = tv.(string)
	}
	if len(token) == 0 {
		for k, v := range m {
			if k == "telegram" {
				if v.Key == "token" {
					cah.Set("telegram_token", v.Value, cache.DefaultExpiration)
					token = v.Value
				}
			}
		}
	}

	return token, err
}
func GetPandoraMessages(ctx context.Context, groupPrefix string) ([]hera_search.TelegramMessage, error) {
	token, err := GetTelegramToken(ctx)
	if err != nil {
		cah.Delete("telegram_token")
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
	if strings.HasPrefix("gpt-4", model) {
		model = "gpt-4"
	}
	if strings.HasPrefix("gpt-3.5", model) {
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

	res := resty_base.GetBaseRestyClient(apiReq.Url, artemis_orchestration_auth.Bearer)
	resp, err := res.R().SetBody(&apiReq.Payload).SetResult(&tc).Post("tokenize")
	if err != nil {
		log.Err(err).Msg("Zeus: GetTokenCountEstimate")
		return 0, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		if err != nil {
			err = fmt.Errorf("Zeus: GetTokenCountEstimate: failed to relay api request: status code %d", resp.StatusCode())
		}
		return 0, err
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
	if len(params.WorkflowInstructions) > 0 {
		systemMessage.Content = params.WorkflowInstructions
	}
	output := hera_search.FormatTgMessagesForAi(msgs)
	tc, err := GetTokenCountEstimate(ctx, output, params.AiModelParams.Model)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	if tc > 10000 {
		return openai.ChatCompletionResponse{}, fmt.Errorf("token count too high: %d", tc)
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
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
		Content: params.WorkflowInstructions,
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	if len(params.WorkflowInstructions) > 0 {
		systemMessage.Content = params.WorkflowInstructions
	}
	output := hera_search.FormatTgMessagesForAi(msgs)
	tc, err := GetTokenCountEstimate(ctx, params.AiModelParams.Model, output)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	if tc > 10000 {
		return openai.ChatCompletionResponse{}, fmt.Errorf("token count too high: %d", tc)
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
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

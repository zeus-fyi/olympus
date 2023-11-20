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
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_access_keygen "github.com/zeus-fyi/olympus/hestia/web/access"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
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

func (h *ZeusAiPlatformActivities) InsertTelegramMessageIfNew(ctx context.Context, ou org_users.OrgUser, msg hera_openai_dbmodels.TelegramMessage) (int, error) {
	tgId, err := hera_openai_dbmodels.InsertNewTgMessages(ctx, ou, msg)
	if err != nil {
		log.Err(err).Interface("msg", msg).Msg("InsertTelegramMessageIfNew: failed")
		return 0, err
	}
	return tgId, nil
}

func (h *ZeusAiPlatformActivities) InsertAiResponse(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

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
func GetPandoraMessages(ctx context.Context, groupPrefix string) ([]hera_openai_dbmodels.TelegramMessage, error) {
	token, err := GetTelegramToken(ctx)
	if err != nil {
		cah.Delete("telegram_token")
		log.Err(err).Msg("Zeus: GetTelegramToken")
		return nil, err
	}

	var msgs []hera_openai_dbmodels.TelegramMessage
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

func GetTokenCountEstimate(ctx context.Context, text string) (int, error) {
	var tc TokenCountsEstimate
	apiReq := &iris_api_requests.ApiProxyRequest{
		Url:             "https://pandora.zeus.fyi",
		PayloadTypeREST: "POST",
		Payload: echo.Map{
			"text": text,
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

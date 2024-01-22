package hera_openai

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

var HeraOpenAI OpenAI

var maxTokensByModel = map[string]int{
	openai.GPT4TurboPreview: 40000,
}

type OpenAI struct {
	*openai.Client
}

func InitHeraOpenAI(bearer string) {
	HeraOpenAI = OpenAI{}
	HeraOpenAI.Client = openai.NewClient(bearer)
}

func InitOrgHeraOpenAI(bearer string) OpenAI {
	return OpenAI{
		openai.NewClient(bearer),
	}
}

type OpenAIParams struct {
	Model                string                    `json:"model"`
	MaxTokens            int                       `json:"maxTokens"`
	Prompt               string                    `json:"prompt"`
	SystemPromptOverride string                    `json:"systemPromptOverride,omitempty"`
	SystemPromptExt      string                    `json:"systemPromptExt,omitempty"`
	FunctionDefinition   openai.FunctionDefinition `json:"functionDefinition,omitempty"`
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

func EnsureValidString(s string) string {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{1,64}$", s)
	if matched {
		return s
	}
	return generateRandomString(10) // You can adjust the length as needed
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func (ai *OpenAI) MakeCodeGenRequestJsonFormattedOutput(ctx context.Context, ou org_users.OrgUser, params OpenAIParams) (openai.ChatCompletionResponse, error) {
	sysPrompt := "Using the specified tool function, analyze the provided input and generate a response that conforms to the expected JSON schema." +
		" The response should precisely reflect the output of the tool function call." +
		" Ensure that your reply is solely in JSON format, containing the results of the function execution, without additional commentary or explanation."
	if params.SystemPromptOverride != "" {
		sysPrompt = params.SystemPromptOverride
	}
	if params.SystemPromptExt != "" {
		sysPrompt += params.SystemPromptExt
	}
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: sysPrompt,
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	params.FunctionDefinition.Name = EnsureValidString(params.FunctionDefinition.Name)
	reqBody := openai.ChatCompletionRequest{
		Model: params.Model,
		Messages: []openai.ChatCompletionMessage{
			systemMessage,
			{
				Role:    openai.ChatMessageRoleFunction,
				Content: params.Prompt,
				Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
			},
		},
		MaxTokens:        0,
		Temperature:      0,
		ResponseFormat:   &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
		FrequencyPenalty: 0,
		LogitBias:        nil,
		Tools: []openai.Tool{{
			Type:     openai.ChatMessageRoleFunction,
			Function: params.FunctionDefinition,
		}},
		ToolChoice: nil,
	}
	if params.MaxTokens > 0 {
		reqBody.MaxTokens = params.MaxTokens
	}
	resp, err := ai.CreateChatCompletion(
		ctx,
		reqBody,
	)
	return resp, err
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
			Model: params.Model,
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

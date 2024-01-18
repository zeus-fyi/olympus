package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

const (
	OpenAiPlatform = "openai"
)

func GetMockingBirdSecrets(ctx context.Context, ou org_users.OrgUser) (*aws_secrets.OAuth2PlatformSecret, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, OpenAiPlatform)
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return nil, err
	}
	return ps, nil
}

func (z *ZeusAiPlatformActivities) SelectTaskDefinition(ctx context.Context, ou org_users.OrgUser, taskID int) ([]artemis_orchestrations.AITaskLibrary, error) {
	tv, err := artemis_orchestrations.SelectTask(ctx, ou, taskID)
	if err != nil {
		log.Err(err).Msg("SelectTaskDefinition: failed to get task definition")
		return nil, err
	}
	return tv, nil
}

func (z *ZeusAiPlatformActivities) AiAnalysisTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, sr []hera_search.SearchResult) (*ChatCompletionQueryResponse, error) {
	cr := openai.ChatCompletionRequest{
		Model:    taskInst.AnalysisModel,
		Messages: []openai.ChatCompletionMessage{},
	}
	content := hera_search.FormatSearchResultsV2(sr)
	if len(content) > 0 {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: taskInst.AnalysisPrompt,
			Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
		}
		cr.Messages = append(cr.Messages, systemMessage)
		// coming from ext
		chatMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
			Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
		}
		cr.Messages = append(cr.Messages, chatMessage)
	} else {
		// else model generated from scratch
		chatMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: taskInst.AnalysisPrompt,
			Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
		}
		cr.Messages = append(cr.Messages, chatMessage)
	}

	if taskInst.AnalysisMaxTokensPerTask > 0 {
		cr.MaxTokens = taskInst.AnalysisMaxTokensPerTask
	}
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "openai")
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return nil, err
	}
	prompt := make(map[string]string)
	prompt["prompt"] = taskInst.AnalysisPrompt
	prompt["content"] = content
	if ps.ApiKey == "" {
		log.Err(err).Msg("AiAnalysisTask: CreateChatCompletion failed, using backup and deleting secret cache for org")
		cres, cerr := hera_openai.HeraOpenAI.CreateChatCompletion(
			ctx, cr,
		)
		if cerr == nil {
			return &ChatCompletionQueryResponse{
				Prompt:         prompt,
				ResponseTaskID: taskInst.AnalysisTaskID,
				Response:       cres,
			}, nil
		} else {
			log.Err(cerr).Msg("AiAnalysisTask: CreateChatCompletion failed")
			return nil, cerr
		}
	}

	oc := hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
	resp, err := oc.CreateChatCompletion(
		ctx, cr,
	)
	if err == nil {
		return &ChatCompletionQueryResponse{
			Prompt:         prompt,
			ResponseTaskID: taskInst.AnalysisTaskID,
			Response:       resp,
		}, nil
	} else {
		log.Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get response using user secrets, clearing cache and trying again")
		aws_secrets.ClearOrgSecretCache(ou)
		ps, err = aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "openai")
		if err != nil || ps == nil || ps.ApiKey == "" {
			if err == nil {
				err = fmt.Errorf("failed to get mockingbird secrets")
			}
			log.Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
			return nil, err
		}
	}
	oc = hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
	resp, err = oc.CreateChatCompletion(
		ctx, cr,
	)
	if err != nil {
		log.Err(err).Msg("AiAnalysisTask: CreateChatCompletion failed")
		return nil, err
	}
	return &ChatCompletionQueryResponse{
		Prompt:         prompt,
		ResponseTaskID: taskInst.AnalysisTaskID,
		Response:       resp,
	}, nil
}

type ChatCompletionQueryResponse struct {
	Prompt                map[string]string             `json:"prompt"`
	Response              openai.ChatCompletionResponse `json:"response"`
	ResponseTaskID        int                           `json:"responseTaskID,omitempty"`
	FilteredMessages      *FilteredMessages             `json:"filteredMessages,omitempty"`
	FilteredSearchResults []hera_search.SearchResult    `json:"filteredSearchResults,omitempty"`

	*twitter.CreateTweetRequest `json:"twitterCreateTweetRequest,omitempty"`
}

func (z *ZeusAiPlatformActivities) AiAggregateTask(ctx context.Context, ou org_users.OrgUser, aggInst artemis_orchestrations.WorkflowTemplateData, dataIn []artemis_orchestrations.AIWorkflowAnalysisResult) (*ChatCompletionQueryResponse, error) {
	content, err := artemis_orchestrations.GenerateContentText(dataIn)
	if err != nil {
		log.Err(err).Msg("AiAggregateTask: GenerateContentText failed")
		return nil, err
	}
	if len(content) <= 0 || aggInst.AggPrompt == nil || aggInst.AggModel == nil || aggInst.AggTaskID == nil || aggInst.AggCycleCount == nil {
		return nil, nil
	}
	if len(*aggInst.AggPrompt) <= 0 {
		return nil, nil
	}
	prompt := make(map[string]string)
	prompt["prompt"] = *aggInst.AggPrompt
	prompt["content"] = content
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: *aggInst.AggPrompt,
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	cr := openai.ChatCompletionRequest{
		Model: *aggInst.AggModel,
		Messages: []openai.ChatCompletionMessage{
			systemMessage,
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
				Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
			},
		},
	}

	if aggInst.AggMaxTokensPerTask == nil {
		aggInst.AggMaxTokensPerTask = aws.Int(0)
	}
	if *aggInst.AggMaxTokensPerTask > 0 {
		cr.MaxTokens = *aggInst.AggMaxTokensPerTask
	}
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, OpenAiPlatform)
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("AiAggregateTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return nil, nil
	}

	if ps.ApiKey == "" {
		log.Err(err).Msg("AiAggregateTask: CreateChatCompletion failed, using backup and deleting secret cache for org")
		cres, cerr := hera_openai.HeraOpenAI.CreateChatCompletion(
			ctx, cr,
		)
		if cerr != nil {
			log.Err(cerr).Msg("AiAggregateTask: CreateChatCompletion failed")
			return nil, cerr
		}
		return &ChatCompletionQueryResponse{
			Prompt:         prompt,
			ResponseTaskID: *aggInst.AggTaskID,
			Response:       cres,
		}, nil
	}

	oc := hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
	resp, err := oc.CreateChatCompletion(
		ctx, cr,
	)
	if err == nil {
		return &ChatCompletionQueryResponse{
			Prompt:         prompt,
			ResponseTaskID: *aggInst.AggTaskID,
			Response:       resp,
		}, nil
	} else {
		log.Err(err).Msg("AiAggregateTask: GetMockingbirdPlatformSecrets: failed to get response using user secrets, clearing cache and trying again")
		aws_secrets.ClearOrgSecretCache(ou)
		ps, err = aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, OpenAiPlatform)
		if err != nil || ps == nil || ps.ApiKey == "" {
			if err == nil {
				err = fmt.Errorf("failed to get mockingbird secrets")
			}
			log.Err(err).Msg("AiAggregateTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
			return nil, err
		}
	}
	oc = hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
	resp, err = oc.CreateChatCompletion(
		ctx, cr,
	)
	if err != nil {
		log.Err(err).Msg("AiAggregateTask: CreateChatCompletion failed")
		return nil, err
	}
	return &ChatCompletionQueryResponse{
		Prompt:         prompt,
		ResponseTaskID: *aggInst.AggTaskID,
		Response:       resp,
	}, nil
}

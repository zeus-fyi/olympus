package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/cvcio/twitter"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hera_twitter "github.com/zeus-fyi/olympus/pkg/hera/twitter"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ZeusAiPlatformActivities struct {
	kronos_helix.ActivityDefinition
}

func NewZeusAiPlatformActivities() ZeusAiPlatformActivities {
	return ZeusAiPlatformActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (z *ZeusAiPlatformActivities) GetActivities() ActivitiesSlice {
	ka := kronos_helix.NewKronosActivities()
	actSlice := []interface{}{z.AiTask, z.SaveAiTaskResponse, z.SendTaskResponseEmail, z.InsertEmailIfNew,
		z.InsertAiResponse, z.InsertTelegramMessageIfNew,
		z.InsertIncomingTweetsFromSearch, z.SearchTwitterUsingQuery, z.SelectTwitterSearchQuery,
		z.SearchRedditNewPostsUsingSubreddit, z.InsertIncomingRedditDataFromSearch, z.SelectRedditSearchQuery,
		z.CreateDiscordJob, z.SelectDiscordSearchQuery, z.InsertIncomingDiscordDataFromSearch,
		z.UpsertAiOrchestration, z.AiAnalysisTask, z.AiRetrievalTask,
		z.AiAggregateTask, z.AiAggregateAnalysisRetrievalTask, z.SaveTaskOutput, z.RecordCompletionResponse,
		z.AiWebRetrievalGetRoutesTask, z.ApiCallRequestTask, z.CreateRedditJob,
		z.SelectActiveSearchIndexerJobs, z.StartIndexingJob, z.CancelRun, z.PlatformIndexerGroupStatusUpdate,
		z.SelectDiscordSearchQueryByGuildChannel, z.CreateJsonOutputModelResponse, z.EvalLookup,
		z.SendResponseToApiForScoresInJson, z.EvalModelScoredJsonOutput,
		z.SendTriggerActionRequestForApproval, z.CreateOrUpdateTriggerActionToExec,
		z.CheckEvalTriggerCondition, z.LookupEvalTriggerConditions,
		z.SocialRedditTask, z.SocialDiscordTask, z.SocialTelegramTask,
		z.SaveTriggerResponseOutput, z.SaveEvalResponseOutput,
		z.SelectTaskDefinition, z.TokenOverflowReduction,
		z.SaveTriggerApiRequestResp, z.SelectRetrievalTask,
		z.SelectTriggerActionToExec, z.SelectTriggerActionApiApprovalWithReqResponses,
		z.CreateOrUpdateTriggerActionApprovalWithApiReq, z.UpdateTriggerActionApproval,
		z.FilterEvalJsonResponses, z.UpdateTaskOutput, z.CreateWsr,
		z.FanOutApiCallRequestTask, z.SaveWorkflowIO, z.SelectWorkflowIO,
		z.SaveCsvTaskOutput,
	}
	return append(actSlice, ka.GetActivities()...)
}

func (z *ZeusAiPlatformActivities) StartIndexingJob(ctx context.Context, sp hera_search.SearchIndexerParams) error {
	switch sp.Platform {
	case redditPlatform:
		ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		err := ZeusAiPlatformWorker.ExecuteAiRedditWorkflow(ctx, ou, sp.SearchGroupName)
		if err != nil {
			log.Err(err).Msg("StartIndexingJob: failed to execute ai reddit workflow")
			return err
		}
	case twitterPlatform:
		ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		err := ZeusAiPlatformWorker.ExecuteAiTwitterWorkflow(ctx, ou, sp.SearchGroupName)
		if err != nil {
			log.Err(err).Msg("StartIndexingJob: failed to execute ai twitter workflow")
			return err
		}
	case telegramPlatform:
		//ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		//err := ZeusAiPlatformWorker.ExecuteAiTelegramWorkflow(ctx, ou, sp.SearchGroupName)
		//if err != nil {
		//	log.Err(err).Msg("StartIndexingJob: failed to execute ai telegram workflow")
		//	return err
		//}
	case discordPlatform:
		ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		err := ZeusAiPlatformWorker.ExecuteAiFetchDataToIngestDiscordWorkflow(ctx, ou, sp.SearchGroupName)
		if err != nil {
			log.Err(err).Msg("StartIndexingJob: failed to execute ai discord workflow")
			return err
		}
	}
	return nil
}

func (z *ZeusAiPlatformActivities) CreateRedditJob(ctx context.Context, ou org_users.OrgUser, subreddit string) error {
	j := RedditJob(subreddit)
	kns := zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "zeus",
		Env:           "production",
	}
	err := zeus.K8Util.DeleteJob(ctx, kns, j.Name)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to delete job")
		return err
	}
	err = zeus.K8Util.DeleteAllPodsLike(ctx, kns, j.Name, nil, nil)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to delete pods")
		return err
	}
	_, err = zeus.K8Util.CreateJob(ctx, kns, &j)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to create job")
		return err
	}
	return err
}

func (z *ZeusAiPlatformActivities) CreateDiscordJob(ctx context.Context, ou org_users.OrgUser, si int, channelID, timeAfter string) error {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "discord")
	if err != nil {
		log.Err(err).Msg("GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		return err
	}
	if ps == nil || ps.ApiKey == "" {
		return fmt.Errorf("GetMockingbirdPlatformSecrets: ps is nil or api key missing")
	}
	hs, err := misc.HashParams([]interface{}{ps.ApiKey})
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to hash params")
		return err
	}
	j := DiscordJob(ou.OrgID, si, ps.ApiKey, hs, channelID, timeAfter)
	kns := zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "zeus",
		Env:           "production"}

	err = zeus.K8Util.DeleteJob(ctx, kns, j.Name)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to delete job")
		return err
	}
	err = zeus.K8Util.DeleteAllPodsLike(ctx, kns, j.Name, nil, nil)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to delete pods")
		return err
	}
	_, err = zeus.K8Util.CreateJob(ctx, kns, &j)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to create job")
		return err
	}
	return err
}

func (z *ZeusAiPlatformActivities) AiTask(ctx context.Context, ou org_users.OrgUser, msg hermes_email_notifications.EmailContents) (openai.ChatCompletionResponse, error) {
	//task := "write a bullet point summary of the email contents and suggest some responses if applicable. write your reply as html formatted\n"
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful bot that reads email contents and provides a bullet point summary and then suggest well thought out responses and that aren't overly formal or stiff in tone and you always write your reply as well formatted html that is easy to read.",
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: Gpt4JsonModel,
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

func (z *ZeusAiPlatformActivities) SendTaskResponseEmail(ctx context.Context, email string, resp openai.ChatCompletionResponse) error {
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

func (z *ZeusAiPlatformActivities) SearchTwitterUsingQuery(ctx context.Context, ou org_users.OrgUser, sp *hera_search.TwitterSearchQuery) ([]*twitter.Tweet, error) {
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, twitterPlatform)
	if err != nil {
		log.Err(err).Msg("SearchTwitterUsingQuery: failed to get mockingbird secrets")
		return nil, err
	}
	if ps == nil || ps.OAuth2Public == "" || ps.OAuth2Secret == "" {
		log.Warn().Interface("ou", ou).Msg("SearchTwitterUsingQuery: ps is empty")
		return nil, fmt.Errorf("SearchTwitterUsingQuery: ps is empty")
	}
	tc, err := hera_twitter.InitTwitterClient(ctx, ps.ConsumerPublic, ps.ConsumerSecret, ps.AccessTokenPublic, ps.AccessTokenSecret)
	if err != nil {
		log.Err(err).Msg("SearchTwitterUsingQuery: failed to init twitter client")
		return nil, err
	}
	tweets, err := tc.GetTweets(ctx, sp.Query, sp.MaxResults, sp.MaxTweetID)
	if err != nil {
		log.Err(err).Msg("SearchTwitterUsingQuery")
		return nil, err
	}
	return tweets, nil
}

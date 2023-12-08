package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/cvcio/twitter"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
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
	}
	return append(actSlice, ka.GetActivities()...)
}

const (
	internalUser = 7138958574876245567
)

func (z *ZeusAiPlatformActivities) CreateDiscordJob(ctx context.Context, si int, channelID, timeAfter string) error {
	authToken, err := read_keys.GetDiscordKey(ctx, internalUser)
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to get discord key")
		return err
	}
	hs, err := misc.HashParams([]interface{}{authToken})
	if err != nil {
		log.Err(err).Msg("CreateDiscordJob: failed to hash params")
		return err
	}
	j := DiscordJob(si, authToken, hs, channelID, timeAfter)
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

func (z *ZeusAiPlatformActivities) SearchRedditNewPostsUsingSubreddit(ctx context.Context, subreddit string, lpo *reddit.ListOptions) ([]*reddit.Post, error) {
	resp, err := hera_reddit.RedditClient.GetNewPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Interface("posts", resp.Posts).Interface("resp", resp.Resp).Msg("SearchRedditNewPostsUsingSubreddit")
		return nil, err
	}
	if resp.Resp.StatusCode >= 400 {
		log.Err(err).Interface("posts", resp.Posts).Interface("resp", resp.Resp).Msg("SearchRedditNewPostsUsingSubreddit")
		return nil, fmt.Errorf("SearchRedditNewPostsUsingSubreddit: resp.StatusCode >= 400")
	}
	return resp.Posts, nil
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

func (z *ZeusAiPlatformActivities) SaveAiTaskResponse(ctx context.Context, ou org_users.OrgUser, resp openai.ChatCompletionResponse, prompt []byte) error {
	err := hera_openai.HeraOpenAI.RecordUIChatRequestUsage(ctx, ou, resp, prompt)
	if err != nil {
		log.Err(err).Msg("SaveAiTaskResponse: RecordUIChatRequestUsage failed")
		return nil
	}
	return nil
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

func (z *ZeusAiPlatformActivities) InsertEmailIfNew(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

func (z *ZeusAiPlatformActivities) InsertTelegramMessageIfNew(ctx context.Context, ou org_users.OrgUser, msg hera_search.TelegramMessage) (int, error) {
	tgId, err := hera_search.InsertNewTgMessages(ctx, ou, msg)
	if err != nil {
		log.Err(err).Interface("msg", msg).Msg("InsertTelegramMessageIfNew: failed")
		return 0, err
	}
	return tgId, nil
}

func (z *ZeusAiPlatformActivities) InsertAiResponse(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

func (z *ZeusAiPlatformActivities) SelectTwitterSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) (*hera_search.TwitterSearchQuery, error) {
	sq, err := hera_search.SelectTwitterSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectTwitterSearchQuery")
		return nil, err
	}
	if sq == nil {
		return nil, fmt.Errorf("SelectTwitterSearchQuery: sq is nil")
	}
	return sq, nil
}

func (z *ZeusAiPlatformActivities) SelectDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) (*hera_search.DiscordSearchResultWrapper, error) {
	sq, err := hera_search.SelectDiscordSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectDiscordSearchQuery")
		return nil, err
	}
	return sq, nil
}

//func (h *ZeusAiPlatformActivities) SelectRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) (*hera_search.TwitterSearchQuery, error) {
//	sq, err := hera_search.SelectRedditSearchQuery(ctx, ou, groupName)
//	if err != nil {
//		log.Err(err).Msg("SelectRedditSearchQuery")
//		return nil, err
//	}
//	if sq == nil {
//		return nil, fmt.Errorf("SelectRedditSearchQuery: sq is nil")
//	}
//	return sq, nil
//}

func (z *ZeusAiPlatformActivities) SearchTwitterUsingQuery(ctx context.Context, sp *hera_search.TwitterSearchQuery) ([]*twitter.Tweet, error) {
	tweets, err := hera_twitter.TwitterClient.GetTweets(ctx, sp.Query, sp.MaxResults, sp.MaxTweetID)
	if err != nil {
		log.Err(err).Msg("SearchTwitterUsingQuery")
		return nil, err
	}
	return tweets, nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingTweetsFromSearch(ctx context.Context, searchID int, tweets []*twitter.Tweet) error {
	_, err := hera_search.InsertIncomingTweets(ctx, searchID, tweets)
	if err != nil {
		log.Err(err).Msg("InsertIncomingTweetsFromSearch")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingRedditDataFromSearch(ctx context.Context, searchID int, redditData []*reddit.Post) error {
	_, err := hera_search.InsertIncomingRedditPosts(ctx, searchID, redditData)
	if err != nil {
		log.Err(err).Msg("InsertIncomingRedditDataFromSearch")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SelectRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) ([]*hera_search.RedditSearchQuery, error) {
	rs, err := hera_search.SelectRedditSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectRedditSearchQuery: activity failed")
		return nil, err
	}
	return rs, nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingDiscordDataFromSearch(ctx context.Context, searchID int, messages hera_discord.ChannelMessages) error {
	err := hera_search.InsertDiscordGuild(ctx, messages.Guild.Id, messages.Guild.Name)
	if err != nil {
		log.Err(err).Msg("InsertIncomingDiscordDataFromSearch")
		return err
	}
	err = hera_search.InsertDiscordChannel(ctx, searchID, messages.Guild.Id, messages.Channel.Id, messages.Channel.CategoryId, messages.Channel.Category, messages.Channel.Name, messages.Channel.Topic)
	if err != nil {
		log.Err(err).Msg("InsertIncomingDiscordDataFromSearch")
		return err
	}
	_, err = hera_search.InsertIncomingDiscordMessages(ctx, searchID, messages)
	if err != nil {
		log.Err(err).Interface("msgs", messages.Messages).Msg("InsertIncomingRedditDataFromSearch")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) UpsertAiOrchestration(ctx context.Context, ou org_users.OrgUser, wfParentID string, wfExecParams artemis_orchestrations.WorkflowExecParams) (int, error) {
	id, err := artemis_orchestrations.UpsertAiOrchestration(ctx, ou, wfParentID, wfExecParams)
	if err != nil {
		log.Err(err).Msg("UpsertAiOrchestration: activity failed")
		return id, err
	}
	return id, nil
}

func (z *ZeusAiPlatformActivities) AiRetrievalTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, window artemis_orchestrations.Window) ([]hera_search.SearchResult, error) {
	if taskInst.RetrievalPlatform == "" || taskInst.RetrievalName == "" || taskInst.RetrievalInstructions == nil {
		return nil, nil
	}
	retInst := artemis_orchestrations.RetrievalItemInstruction{}
	jerr := json.Unmarshal(taskInst.RetrievalInstructions, &retInst)
	if jerr != nil {
		log.Err(jerr).Msg("AiRetrievalTask: failed to unmarshal")
		return nil, jerr
	}
	sp := hera_search.AiSearchParams{
		Retrieval: artemis_orchestrations.RetrievalItem{
			RetrievalItemInstruction: retInst,
		},
		Window: window,
	}

	var resp []hera_search.SearchResult
	var err error
	switch taskInst.RetrievalPlatform {
	case "twitter":
		resp, err = hera_search.SearchTwitter(ctx, ou, sp)
	case "reddit":
		resp, err = hera_search.SearchReddit(ctx, ou, sp)
	case "discord":
		resp, err = hera_search.SearchDiscord(ctx, ou, sp)
	case "telegram":
		resp, err = hera_search.SearchTelegram(ctx, ou, sp)
	case "web":
		// todo
	}
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	return resp, nil
}

func (z *ZeusAiPlatformActivities) AiAnalysisTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, sr []hera_search.SearchResult) (openai.ChatCompletionResponse, error) {
	content := hera_search.FormatSearchResultsV2(sr)
	if len(content) <= 0 {
		return openai.ChatCompletionResponse{}, nil
	}
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: taskInst.AnalysisPrompt,
		Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
	}
	cr := openai.ChatCompletionRequest{
		Model: taskInst.AnalysisModel,
		Messages: []openai.ChatCompletionMessage{
			systemMessage,
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
				Name:    fmt.Sprintf("%d-%d", ou.OrgID, ou.UserID),
			},
		},
	}
	if taskInst.AnalysisMaxTokensPerTask > 0 {
		cr.MaxTokens = taskInst.AnalysisMaxTokensPerTask
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx, cr,
	)
	if err != nil {
		log.Err(err).Msg("AiAnalysisTask: CreateChatCompletion failed")
		return resp, err
	}
	return resp, err
}

func (z *ZeusAiPlatformActivities) AiAggregateTask(ctx context.Context, ou org_users.OrgUser, aggInst artemis_orchestrations.WorkflowTemplateData, dataIn []artemis_orchestrations.AIWorkflowAnalysisResult) (openai.ChatCompletionResponse, error) {
	content := artemis_orchestrations.GenerateContentText(dataIn)
	if len(content) <= 0 || aggInst.AggPrompt == nil || aggInst.AggModel == nil || aggInst.AggTaskID == nil || aggInst.AggCycleCount == nil {
		return openai.ChatCompletionResponse{}, nil
	}
	if len(*aggInst.AggPrompt) <= 0 {
		return openai.ChatCompletionResponse{}, nil
	}
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
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		ctx, cr,
	)
	if err != nil {
		log.Err(err).Msg("AiAggregateTask: CreateChatCompletion failed")
		return resp, err
	}
	return resp, err
}

func (z *ZeusAiPlatformActivities) RecordCompletionResponse(ctx context.Context, ou org_users.OrgUser, resp openai.ChatCompletionResponse, prompt []byte) (int, error) {
	rid, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, ou, resp, prompt)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformActivities: RecordCompletionResponse: failed")
		return rid, err
	}
	return rid, nil
}

func (z *ZeusAiPlatformActivities) AiAggregateAnalysisRetrievalTask(ctx context.Context, window artemis_orchestrations.Window, ojIDs, sourceTaskIds []int) ([]artemis_orchestrations.AIWorkflowAnalysisResult, error) {
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, window, ojIDs, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: failed")
		return nil, err
	}

	return results, nil
}

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr artemis_orchestrations.AIWorkflowAnalysisResult) error {
	respID, err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return err
	}
	return nil
}

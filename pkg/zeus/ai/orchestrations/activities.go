package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/cvcio/twitter"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hera_twitter "github.com/zeus-fyi/olympus/pkg/hera/twitter"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
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
		z.AiWebRetrievalGetRoutesTask, z.AiWebRetrievalTask, z.CreateRedditJob,
		z.SelectActiveSearchIndexerJobs, z.StartIndexingJob, z.CancelRun, z.PlatformIndexerGroupStatusUpdate,
		z.SelectDiscordSearchQueryByGuildChannel, z.CreateJsonOutputModelResponse, z.EvalLookup,
		z.SendResponseToApiForScoresInJson, z.EvalModelScoredJsonOutput, z.SaveEvalMetricResults,
		z.SendTriggerActionRequestForApproval, z.CreateOrUpdateTriggerActionToExec,
		z.CheckEvalTriggerCondition, z.LookupEvalTriggerConditions,
		z.SocialTweetTask, z.SocialRedditTask, z.SocialDiscordTask,
	}
	return append(actSlice, ka.GetActivities()...)
}

func (z *ZeusAiPlatformActivities) StartIndexingJob(ctx context.Context, sp hera_search.SearchIndexerParams) error {
	switch sp.Platform {
	case "reddit":
		ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		err := ZeusAiPlatformWorker.ExecuteAiRedditWorkflow(ctx, ou, sp.SearchGroupName)
		if err != nil {
			log.Err(err).Msg("StartIndexingJob: failed to execute ai reddit workflow")
			return err
		}
	case "twitter":
		ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		err := ZeusAiPlatformWorker.ExecuteAiTwitterWorkflow(ctx, ou, sp.SearchGroupName)
		if err != nil {
			log.Err(err).Msg("StartIndexingJob: failed to execute ai twitter workflow")
			return err
		}
	case "telegram":
		//ou := org_users.NewOrgUserWithID(sp.OrgID, 0)
		//err := ZeusAiPlatformWorker.ExecuteAiTelegramWorkflow(ctx, ou, sp.SearchGroupName)
		//if err != nil {
		//	log.Err(err).Msg("StartIndexingJob: failed to execute ai telegram workflow")
		//	return err
		//}
	case "discord":
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
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "twitter")
	if err != nil {
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to get mockingbird secrets")
		return nil, err
	}
	if ps == nil || ps.OAuth2Public == "" || ps.OAuth2Secret == "" {
		log.Warn().Interface("ou", ou).Msg("SearchTwitterUsingQuery: ps is empty")
		return nil, fmt.Errorf("SearchTwitterUsingQuery: ps is empty")
	}
	tc, err := hera_twitter.InitOrgTwitterClient(ctx, ps.OAuth2Public, ps.OAuth2Secret)
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

func (z *ZeusAiPlatformActivities) AiWebRetrievalGetRoutesTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData) ([]iris_models.RouteInfo, error) {
	retInst := artemis_orchestrations.RetrievalItemInstruction{}
	jerr := json.Unmarshal(taskInst.RetrievalInstructions, &retInst)
	if jerr != nil {
		log.Err(jerr).Msg("AiRetrievalTask: failed to unmarshal")
		return nil, jerr
	}
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, jerr
	}
	ogr, rerr := iris_models.SelectOrgGroupRoutes(ctx, ou.OrgID, *retInst.WebFilters.RoutingGroup)
	if rerr != nil {
		log.Err(rerr).Msg("AiRetrievalTask: failed to select org routes")
		return nil, rerr
	}
	return ogr, nil

}

func (z *ZeusAiPlatformActivities) AiWebRetrievalTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, r iris_models.RouteInfo) (*hera_search.SearchResult, error) {
	retInst := artemis_orchestrations.RetrievalItemInstruction{}
	jerr := json.Unmarshal(taskInst.RetrievalInstructions, &retInst)
	if jerr != nil {
		log.Err(jerr).Msg("AiRetrievalTask: failed to unmarshal")
		return nil, jerr
	}
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, jerr
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	req := &iris_api_requests.ApiProxyRequest{
		Url:             r.RoutePath,
		PayloadTypeREST: "GET",
		Timeout:         1 * time.Minute,
		StatusCode:      http.StatusOK,
	}
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, fmt.Sprintf("web-%s", *retInst.WebFilters.RoutingGroup))
	if err == nil && ps != nil {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("SearchRedditNewPostsUsingSubreddit: failed to get mockingbird secrets")
		if ps != nil && ps.ApiKey != "" {
			req.Bearer = ps.ApiKey
		}
	} else {
		err = nil
	}

	rr, rrerr := rw.ExtLoadBalancerRequest(ctx, req)
	if rrerr != nil {
		log.Err(rrerr).Msg("AiRetrievalTask: failed to request")
		return nil, rrerr
	}
	wr := hera_search.WebResponse{
		Body:       rr.Response,
		RawMessage: rr.RawResponse,
	}
	value := ""
	if wr.Body != nil {
		b, jer := json.Marshal(wr.Body)
		if jer != nil {
			log.Err(jer).Msg("AiRetrievalTask: failed to marshal")
			return nil, jer
		}
		value = fmt.Sprintf("%s", b)
	}
	if wr.RawMessage != nil && wr.Body == nil {
		value = fmt.Sprintf("%s", wr.RawMessage)
	}

	sres := &hera_search.SearchResult{
		Source:      rr.Url,
		Value:       value,
		Group:       aws.StringValue(retInst.WebFilters.RoutingGroup),
		WebResponse: wr,
	}

	return sres, nil
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
	default:
		return nil, nil
	}
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	return resp, nil
}

func (z *ZeusAiPlatformActivities) RecordCompletionResponse(ctx context.Context, ou org_users.OrgUser, resp *ChatCompletionQueryResponse) (int, error) {
	if resp == nil {
		return 0, nil
	}
	b, err := json.Marshal(resp.Prompt)
	if err != nil {
		log.Err(err).Msg("RecordCompletionResponse: failed")
		return 0, err
	}
	rid, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, ou, resp.Response, b)
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

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr artemis_orchestrations.AIWorkflowAnalysisResult, dataIn []artemis_orchestrations.AIWorkflowAnalysisResult) error {
	md, err := json.Marshal(dataIn)
	if err != nil {
		log.Err(err).Interface("dataIn", dataIn).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return err
	}
	wr.Metadata = md
	respID, err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return err
	}
	return nil
}

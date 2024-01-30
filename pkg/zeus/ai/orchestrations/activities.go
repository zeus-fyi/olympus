package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/cvcio/twitter"
	"github.com/labstack/echo/v4"
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
		z.AiWebRetrievalGetRoutesTask, z.ApiCallRequestTask, z.CreateRedditJob,
		z.SelectActiveSearchIndexerJobs, z.StartIndexingJob, z.CancelRun, z.PlatformIndexerGroupStatusUpdate,
		z.SelectDiscordSearchQueryByGuildChannel, z.CreateJsonOutputModelResponse, z.EvalLookup,
		z.SendResponseToApiForScoresInJson, z.EvalModelScoredJsonOutput, z.SaveEvalMetricResults,
		z.SendTriggerActionRequestForApproval, z.CreateOrUpdateTriggerActionToExec,
		z.CheckEvalTriggerCondition, z.LookupEvalTriggerConditions,
		z.SocialTweetTask, z.SocialRedditTask, z.SocialDiscordTask, z.SocialTelegramTask,
		z.EvalFormatForApi, z.SaveTriggerResponseOutput, z.SaveEvalResponseOutput,
		z.SelectTaskDefinition, z.ExtractTweets, z.TokenOverflowReduction,
		z.AnalyzeEngagementTweets, z.SaveTriggerApiRequestResp, z.SelectRetrievalTask,
		z.SelectTriggerActionToExec, z.SelectTriggerActionApiApprovalWithReqResponses,
		z.CreateOrUpdateTriggerActionApprovalWithApiReq,
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

func (z *ZeusAiPlatformActivities) AiWebRetrievalGetRoutesTask(ctx context.Context, ou org_users.OrgUser, retrieval artemis_orchestrations.RetrievalItem) ([]iris_models.RouteInfo, error) {
	if retrieval.WebFilters == nil || retrieval.WebFilters.RoutingGroup == nil || len(*retrieval.WebFilters.RoutingGroup) <= 0 {
		return nil, nil
	}
	ogr, rerr := iris_models.SelectOrgGroupRoutes(ctx, ou.OrgID, *retrieval.WebFilters.RoutingGroup)
	if rerr != nil {
		log.Err(rerr).Msg("AiRetrievalTask: failed to select org routes")
		return nil, rerr
	}
	return ogr, nil

}

type RouteTask struct {
	Ou        org_users.OrgUser                    `json:"orgUser"`
	Retrieval artemis_orchestrations.RetrievalItem `json:"retrieval"`
	RouteInfo iris_models.RouteInfo                `json:"routeInfo"`
	Payload   echo.Map                             `json:"payload"`
}

func (z *ZeusAiPlatformActivities) ApiCallRequestTask(ctx context.Context, r RouteTask) (*hera_search.SearchResult, error) {
	retInst := artemis_orchestrations.RetrievalItemInstruction{}
	jerr := json.Unmarshal(r.Retrieval.RetrievalItemInstruction.Instructions.Bytes, &retInst)
	if jerr != nil {
		log.Err(jerr).Msg("AiRetrievalTask: failed to unmarshal")
		return nil, jerr
	}
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, jerr
	}
	restMethod := http.MethodGet
	if retInst.WebFilters.EndpointREST != nil {
		restMethod = *retInst.WebFilters.EndpointREST
	}
	var routeExt string
	if retInst.WebFilters.EndpointREST != nil {
		routeExt = *retInst.WebFilters.EndpointRoutePath
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	req := &iris_api_requests.ApiProxyRequest{
		Url:             r.RouteInfo.RoutePath,
		OrgID:           r.Ou.OrgID,
		UserID:          r.Ou.UserID,
		ExtRoutePath:    routeExt,
		Payload:         r.Payload,
		PayloadTypeREST: restMethod,
	}
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, r.Ou, fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup))
	if err == nil && ps != nil {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		log.Err(err).Msg("AiWebRetrievalTask: failed to get mockingbird secrets")
		if ps != nil && ps.ApiKey != "" {
			req.Bearer = ps.ApiKey
		}
	} else {
		err = nil
	}

	rr, rrerr := rw.ExtLoadBalancerRequest(ctx, req)
	if rrerr != nil {
		log.Err(rrerr).Msg("AiWebRetrievalTask: failed to request")
		return nil, rrerr
	}
	wr := hera_search.WebResponse{
		WebFilters: retInst.WebFilters,
		Body:       rr.Response,
		RawMessage: rr.RawResponse,
	}
	value := ""
	if wr.Body != nil {
		b, jer := json.Marshal(wr.Body)
		if jer != nil {
			log.Err(jer).Msg("AiWebRetrievalTask: failed to marshal")
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

func (z *ZeusAiPlatformActivities) AiRetrievalTask(ctx context.Context, ou org_users.OrgUser,
	retrieval artemis_orchestrations.RetrievalItem, window artemis_orchestrations.Window) (*hera_search.SearchResultGroup, error) {
	if retrieval.RetrievalPlatform == "" || retrieval.RetrievalName == "" {
		return nil, nil
	}
	sg := &hera_search.SearchResultGroup{
		PlatformName: retrieval.RetrievalPlatform,
	}
	sp := hera_search.AiSearchParams{
		Retrieval: artemis_orchestrations.RetrievalItem{
			RetrievalItemInstruction: retrieval.RetrievalItemInstruction,
		},
		Window: window,
	}
	if ou.OrgID == 7138983863666903883 && retrieval.RetrievalName == "twitter-test" {
		aiSp := hera_search.AiSearchParams{
			TimeRange: "30 days",
		}
		hera_search.TimeRangeStringToWindow(&aiSp)
		resp, err := hera_search.SearchTwitter(ctx, ou, aiSp)
		if err != nil {
			log.Err(err).Msg("AiRetrievalTask: failed")
			return nil, err
		}
		if len(resp) > 50 {
			resp = resp[:50]
		}
		sg.SearchResults = resp
		return sg, nil
	}

	var resp []hera_search.SearchResult
	var err error
	switch retrieval.RetrievalPlatform {
	case twitterPlatform:
		resp, err = hera_search.SearchTwitter(ctx, ou, sp)
	case redditPlatform:
		resp, err = hera_search.SearchReddit(ctx, ou, sp)
	case discordPlatform:
		resp, err = hera_search.SearchDiscord(ctx, ou, sp)
	case telegramPlatform:
		resp, err = hera_search.SearchTelegram(ctx, ou, sp)
	default:
		return nil, nil
	}
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	sg.SearchResults = resp
	return sg, nil
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

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr *artemis_orchestrations.AIWorkflowAnalysisResult, dataIn any) (int, error) {
	if wr == nil {
		return 0, nil
	}
	md, err := json.Marshal(dataIn)
	if err != nil {
		log.Err(err).Interface("dataIn", dataIn).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	wr.Metadata = md
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	return wr.WorkflowResultID, nil
}

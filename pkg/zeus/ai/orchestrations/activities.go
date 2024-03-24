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
		z.SendResponseToApiForScoresInJson, z.EvalModelScoredJsonOutput,
		z.SendTriggerActionRequestForApproval, z.CreateOrUpdateTriggerActionToExec,
		z.CheckEvalTriggerCondition, z.LookupEvalTriggerConditions,
		z.SocialRedditTask, z.SocialDiscordTask, z.SocialTelegramTask,
		z.SaveTriggerResponseOutput, z.SaveEvalResponseOutput,
		z.SelectTaskDefinition, z.TokenOverflowReduction,
		z.SaveTriggerApiRequestResp, z.SelectRetrievalTask,
		z.SelectTriggerActionToExec, z.SelectTriggerActionApiApprovalWithReqResponses,
		z.CreateOrUpdateTriggerActionApprovalWithApiReq, z.UpdateTriggerActionApproval,
		z.FilterEvalJsonResponses, z.UpdateTaskOutput,
		z.SelectWorkflowIO, z.SaveWorkflowIO,
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
	Payloads  []echo.Map                           `json:"payloads"`
	Headers   http.Header                          `json:"headers"`
}

func FixRegexInput(input string) string {
	if len(input) > 0 {
		// Check if the first character is a backtick and replace it with a double quote
		if input[0] == '`' {
			input = "\"" + input[1:]
		}
		// Check if the last character is a backtick and replace it with a double quote
		if input[len(input)-1] == '`' {
			input = input[:len(input)-1] + "\""
		}
	}
	return input
}
func (z *ZeusAiPlatformActivities) ApiCallRequestTask(ctx context.Context, r RouteTask, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	retInst := r.Retrieval
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, nil
	}
	restMethod := http.MethodGet
	if retInst.WebFilters.EndpointREST != nil {
		restMethod = strings.ToLower(*retInst.WebFilters.EndpointREST)
		switch restMethod {
		case "post", "POST":
			restMethod = http.MethodPost
		case "put", "PUT":
			restMethod = http.MethodPut
		case "delete":
			restMethod = http.MethodDelete
		case "patch":
			restMethod = http.MethodPatch
		case "get":
			restMethod = http.MethodGet
		default:
			log.Info().Str("restMethod", restMethod).Msg("ApiCallRequestTask: rest method")
		}
	}
	var routeExt string
	var orgRouteExt string
	if retInst.WebFilters.EndpointREST != nil {
		routeExt = *retInst.WebFilters.EndpointRoutePath
	}
	var extractedQpsVals []string
	if r.Payload != nil && retInst.WebFilters.RegexPatterns != nil {
		rp, qps, err := ReplaceAndPassParams(routeExt, r.Payload)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to replace route path params")
			return nil, err
		}

		log.Info().Interface("qps", qps).Msg("ApiCallRequestTask: qps")
		extractedQpsVals = qps
		orgRouteExt = routeExt
		routeExt = rp
		if len(r.Payload) == 0 {
			r.Payload = nil
		}
	}

	var sec []int
	if retInst.WebFilters.DontRetryStatusCodes != nil {
		sec = retInst.WebFilters.DontRetryStatusCodes
	}
	secretNameRefApi := fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)
	var regexPatterns []string
	for _, rgp := range retInst.WebFilters.RegexPatterns {
		regexPatterns = append(regexPatterns, FixRegexInput(rgp))
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	req := &iris_api_requests.ApiProxyRequest{
		Url:                    r.RouteInfo.RoutePath,
		OrgID:                  r.Ou.OrgID,
		UserID:                 r.Ou.UserID,
		ExtRoutePath:           routeExt,
		Payload:                r.Payload,
		Payloads:               r.Payloads,
		PayloadTypeREST:        restMethod,
		RequestHeaders:         r.Headers,
		RegexFilters:           regexPatterns,
		SkipErrorOnStatusCodes: sec,
		SecretNameRef:          secretNameRefApi,
	}
	rr, rrerr := rw.ExtLoadBalancerRequest(ctx, req)
	if rrerr != nil {
		if req.StatusCode == 401 {
			// clear the cache
			log.Warn().Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Int("statusCode", req.StatusCode).Msg("ApiCallRequestTask: clearing org secret cache")
			aws_secrets.ClearOrgSecretCache(r.Ou)
		}
		log.Err(rrerr).Interface("payload", r.Payload).Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Msg("ApiCallRequestTask: failed to get response")
		return nil, rrerr
	}
	rr.ExtRoutePath = orgRouteExt
	req.ExtRoutePath = orgRouteExt
	wr := hera_search.WebResponse{
		WebFilters: retInst.WebFilters,
		Body:       rr.Response,
		RawMessage: rr.RawResponse,
	}

	value := ""
	if wr.Body != nil && wr.RawMessage == nil {
		b, jer := json.Marshal(wr.Body)
		if jer != nil {
			log.Err(jer).Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Msg("ApiCallRequestTask: failed to get response")
			return nil, jer
		}
		value = fmt.Sprintf("%s", b)
	} else if wr.RawMessage != nil && len(req.RegexFilters) > 0 {
		value = fmt.Sprintf("%s", wr.RawMessage)
		wr.RegexFilteredBody = value
	} else if wr.Body != nil && wr.RawMessage != nil {
		// todo & redundant
		value = fmt.Sprintf("%s", wr.RawMessage)
	}

	sres := hera_search.SearchResult{
		Source:      rr.Url,
		Value:       value,
		QueryParams: extractedQpsVals,
		Group:       aws.StringValue(retInst.WebFilters.RoutingGroup),
		WebResponse: wr,
	}
	sg := &hera_search.SearchResultGroup{
		PlatformName: cp.Tc.Retrieval.RetrievalPlatform,
		Window:       cp.Window,
	}
	sg.ApiResponseResults = []hera_search.SearchResult{sres}
	if req.RegexFilters != nil && len(req.RegexFilters) > 0 {
		cp.Tc.RegexSearchResults = append(cp.Tc.RegexSearchResults, sres)
	} else {
		cp.Tc.ApiResponseResults = append(cp.Tc.ApiResponseResults, sres)
	}
	sg.SourceTaskID = cp.Tc.TaskID
	if cp.Wsr.InputID == 0 {
		wio := WorkflowStageIO{
			WorkflowStageReference: cp.Wsr,
			WorkflowStageInfo: WorkflowStageInfo{
				PromptReduction: &PromptReduction{
					MarginBuffer:          cp.Tc.MarginBuffer,
					Model:                 cp.Tc.Model,
					TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
					PromptReductionSearchResults: &PromptReductionSearchResults{
						InPromptBody:  cp.Tc.Prompt,
						InSearchGroup: sg,
					},
				},
			},
		}
		wid, err := sws(ctx, &wio)
		if err != nil {
			log.Err(err).Msg("AiRetrievalTask: failed")
			return nil, err
		}
		cp.Wsr.InputID = wid.InputID
	} else {
		wio, werr := gws(ctx, cp.Wsr.InputID)
		if werr != nil {
			log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
			return nil, werr
		}
		if wio.WorkflowStageInfo.PromptReduction == nil {
			wio.WorkflowStageInfo.PromptReduction = &PromptReduction{
				MarginBuffer:          cp.Tc.MarginBuffer,
				Model:                 cp.Tc.Model,
				TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
				PromptReductionSearchResults: &PromptReductionSearchResults{
					InPromptBody:  cp.Tc.Prompt,
					InSearchGroup: sg,
				},
			}
		} else if wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults == nil || wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup == nil {
			wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults = &PromptReductionSearchResults{
				InPromptBody:  cp.Tc.Prompt,
				InSearchGroup: sg,
			}
		} else {
			if wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults == nil {
				wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults = make([]hera_search.SearchResult, 0)
			}
			wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults = append(wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults, sres)
		}
		wo, err := sws(ctx, &wio)
		if err != nil {
			log.Err(err).Msg("TokenOverflowReduction: failed to update workflow io")
			return nil, err
		}
		cp.Wsr.InputID = wo.InputID
	}
	return cp, nil
}

func (z *ZeusAiPlatformActivities) AiRetrievalTask(ctx context.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	retrieval := cp.Tc.Retrieval
	ou := cp.Ou
	window := cp.Window

	if retrieval.RetrievalPlatform == "" || retrieval.RetrievalName == "" {
		return nil, nil
	}
	sg := &hera_search.SearchResultGroup{
		PlatformName: retrieval.RetrievalPlatform,
		Window:       window,
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
		sg.SourceTaskID = cp.Tc.TaskID
		wio := WorkflowStageIO{
			WorkflowStageReference: cp.Wsr,
			WorkflowStageInfo: WorkflowStageInfo{
				PromptReduction: &PromptReduction{
					MarginBuffer:          cp.Tc.MarginBuffer,
					Model:                 cp.Tc.Model,
					TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
					PromptReductionSearchResults: &PromptReductionSearchResults{
						InPromptBody:  cp.Tc.Prompt,
						InSearchGroup: sg,
					},
				},
			},
		}
		wid, err := sws(ctx, &wio)
		if err != nil {
			log.Err(err).Msg("AiRetrievalTask: failed")
			return nil, err
		}
		cp.Wsr.InputID = wid.InputID
		return cp, nil
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
	sg.SourceTaskID = cp.Tc.TaskID
	wio := WorkflowStageIO{
		WorkflowStageReference: cp.Wsr,
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				MarginBuffer:          cp.Tc.MarginBuffer,
				Model:                 cp.Tc.Model,
				TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
				PromptReductionSearchResults: &PromptReductionSearchResults{
					InPromptBody:  cp.Tc.Prompt,
					InSearchGroup: sg,
				},
			},
		},
	}
	wid, err := sws(ctx, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	cp.Wsr.InputID = wid.InputID
	return cp, nil
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

type AggRetResp struct {
	AIWorkflowAnalysisResultSlice []artemis_orchestrations.AIWorkflowAnalysisResult
	InputDataAnalysisToAggSlice   []InputDataAnalysisToAgg
}

func (z *ZeusAiPlatformActivities) AiAggregateAnalysisRetrievalTask(ctx context.Context, cp *MbChildSubProcessParams, sourceTaskIds []int) (*MbChildSubProcessParams, error) {
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, cp.Window, []int{cp.Oj.OrchestrationID}, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: failed")
		return nil, err
	}
	var resp []InputDataAnalysisToAgg
	for _, r := range results {
		b, berr := json.Marshal(r.Metadata)
		if berr != nil {
			log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		tmp := InputDataAnalysisToAgg{}
		jerr := json.Unmarshal(b, &tmp)
		if jerr != nil {
			log.Err(jerr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		resp = append(resp, tmp)
	}
	wio := WorkflowStageIO{
		WorkflowStageReference: cp.Wsr,
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				MarginBuffer:              cp.Tc.MarginBuffer,
				Model:                     cp.Tc.Model,
				TokenOverflowStrategy:     cp.Tc.TokenOverflowStrategy,
				DataInAnalysisAggregation: resp,
			},
		},
	}
	wid, err := sws(ctx, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	cp.Wsr.InputID = wid.InputID
	return cp, nil
}

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr *artemis_orchestrations.AIWorkflowAnalysisResult, cp *MbChildSubProcessParams, dataIn InputDataAnalysisToAgg) (int, error) {
	if cp == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp is nil")
	}
	wio, werr := gws(ctx, cp.Wsr.InputID)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
	}
	if wio.PromptReduction != nil && wio.PromptReduction.PromptReductionSearchResults != nil && wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups != nil && len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) > 0 {
		if cp.Wsr.ChunkOffset < len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) {
			dataIn.SearchResultGroup = wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups[cp.Wsr.ChunkOffset]
		}
	}
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

	if dataIn.TextInput != nil {
		wio.PromptTextFromTextStage += *dataIn.TextInput
		_, eerr := sws(ctx, &wio)
		if eerr != nil {
			log.Err(err).Msg("SaveTaskOutput: failed")
			return -1, eerr
		}
	}
	return wr.WorkflowResultID, nil
}

/*
	var zs string
	if dataIn.ChatCompletionQueryResponse != nil && dataIn.ChatCompletionQueryResponse.Response.Choices != nil {
		for _, c := range dataIn.ChatCompletionQueryResponse.Response.Choices {
			zs += c.Message.Content
		}
		if wio.PromptReduction == nil {
			wio.PromptReduction = &PromptReduction{}
		}
		if wio.PromptReduction.PromptReductionText == nil {
			wio.PromptReduction.PromptReductionText = &PromptReductionText{}
		}
		wio.PromptReduction.PromptReductionText.InPromptBody += zs
		_, err = sws(ctx, &wio)
		if err != nil {
			log.Err(err).Msg("AiRetrievalTask: failed")
			return 0, err
		}
	}
*/

// UpdateTaskOutput updates the task output, but it only intended for json output results
func (z *ZeusAiPlatformActivities) UpdateTaskOutput(ctx context.Context, cp *MbChildSubProcessParams) ([]artemis_orchestrations.JsonSchemaDefinition, error) {
	if cp == nil || len(cp.Tc.JsonResponseResults) <= 0 {
		return nil, nil
	}
	var skipAnalysis bool
	jro := FilterPassingEvalPassingResponses(cp.Tc.JsonResponseResults)
	var md []byte
	var err error
	var filteredJsonResponses []artemis_orchestrations.JsonSchemaDefinition
	var infoJsonResponses []artemis_orchestrations.JsonSchemaDefinition
	for evalState, v := range jro {
		switch evalState {
		case filterState:
			filteredJsonResponses = v.Passed
			log.Info().Interface("v.Passed", len(v.Passed)).Interface("v.Failed", len(v.Failed)).Msg("UpdateTaskOutput: filterState")
		case infoState:
			if len(v.Failed) > 0 {
				skipAnalysis = true
			} else {
				infoJsonResponses = v.Passed
				log.Info().Interface("v.Passed", len(v.Passed)).Msg("UpdateTaskOutput: infoState")
			}
		case errorState:
			// TODO: stop workflow?
		}
	}
	var sg *hera_search.SearchResultGroup
	wio, werr := gws(ctx, cp.Wsr.InputID)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return nil, werr
	}
	if wio.PromptReduction != nil && wio.PromptReduction.PromptReductionSearchResults != nil && wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups != nil && len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) > 0 {
		if cp.Wsr.ChunkOffset < len(wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups) {
			sg = wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups[cp.Wsr.ChunkOffset]
		}
	}
	var res []artemis_orchestrations.JsonSchemaDefinition
	if len(filteredJsonResponses) <= 0 && len(infoJsonResponses) <= 0 {
		skipAnalysis = true
		md, err = json.Marshal(jro)
		if err != nil {
			log.Err(err).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	} else if len(filteredJsonResponses) > 0 {
		res = filteredJsonResponses
		tmp := InputDataAnalysisToAgg{
			SearchResultGroup: sg,
			ChatCompletionQueryResponse: &ChatCompletionQueryResponse{
				JsonResponseResults: res,
				RegexSearchResults:  cp.Tc.RegexSearchResults,
			},
		}
		md, err = json.Marshal(tmp)
		if err != nil {
			log.Err(err).Interface("infoJsonResponses", infoJsonResponses).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	} else {
		res = infoJsonResponses
		tmp := InputDataAnalysisToAgg{
			SearchResultGroup: sg,
			ChatCompletionQueryResponse: &ChatCompletionQueryResponse{
				JsonResponseResults: res,
				RegexSearchResults:  cp.Tc.RegexSearchResults,
			},
		}
		md, err = json.Marshal(tmp)
		if err != nil {
			log.Err(err).Interface("infoJsonResponses", infoJsonResponses).Interface("jro", jro).Msg("UpdateTaskOutput: failed")
			return nil, err
		}
	}
	if res != nil && sg != nil && sg.SearchResults != nil {
		seen := make(map[int]bool)
		for _, jr := range res {
			for _, fv := range jr.Fields {
				if fv.FieldName == "msg_id" && fv.IsValidated && fv.NumberValue != nil && *fv.NumberValue > 0 {
					seen[int(*fv.NumberValue)] = true
				}
				if fv.FieldName == "msg_id" && fv.IsValidated && fv.IntegerValue != nil && *fv.IntegerValue > 0 {
					seen[*fv.IntegerValue] = true
				}
			}
		}
		sg.FilteredSearchResults = []hera_search.SearchResult{}
		for _, sr := range sg.SearchResults {
			_, ok := seen[sr.UnixTimestamp]
			if ok {
				sg.FilteredSearchResults = append(sg.FilteredSearchResults, sr)
			}
		}
	}
	wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
		WorkflowResultID:      cp.Tc.WorkflowResultID,
		ResponseID:            cp.Tc.ResponseID,
		OrchestrationID:       cp.Oj.OrchestrationID,
		SourceTaskID:          cp.Tc.TaskID,
		IterationCount:        cp.Wsr.IterationCount,
		ChunkOffset:           cp.Wsr.ChunkOffset,
		RunningCycleNumber:    cp.Wsr.RunCycle,
		SearchWindowUnixStart: cp.Window.UnixStartTime,
		SearchWindowUnixEnd:   cp.Window.UnixEndTime,
		Metadata:              md,
		SkipAnalysis:          skipAnalysis,
	}
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("filteredJsonResponses", filteredJsonResponses).Interface("jro", jro).Interface("wr", wr).Msg("UpdateTaskOutput: failed")
		return nil, err
	}
	return res, nil
}

package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	Headers   http.Header                          `json:"headers"`
	Qps       []string                             `json:"qps"`
	PrevRoute string                               `json:"prevRoute,omitempty"`
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
		_, err = s3ws(ctx, cp, &wio)
		if err != nil {
			log.Err(err).Msg("AiRetrievalTask: failed")
			return nil, err
		}
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
	_, err = s3ws(ctx, cp, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
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
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
		return nil, err
	}
	var resp []InputDataAnalysisToAgg
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
		for _, vi := range cp.WfExecParams.WorkflowTasks {
			if vi.AnalysisTaskName != "" {
				tn := cp.Tc.TaskName
				cp.Tc.TaskName = vi.AnalysisTaskName
				wso, werr := gs3wfs(ctx, cp)
				if werr != nil {
					log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
					return nil, err
				}
				resp = append(resp, wso.InputDataAnalysisToAgg)
				cp.Tc.TaskName = tn
			}
		}
	} else {
		for _, r := range results {
			b, berr := json.Marshal(r.Metadata)
			if berr != nil {
				log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
				continue
			}
			if b == nil || string(b) == "null" {
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
	}
	log.Info().Interface("len(results)", len(results)).Msg("AiAggregateAnalysisRetrievalTask")
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
	_, err = s3ws(ctx, cp, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	return cp, nil
}

func (z *ZeusAiPlatformActivities) SaveTaskOutput(ctx context.Context, wr *artemis_orchestrations.AIWorkflowAnalysisResult, cp *MbChildSubProcessParams, dataIn InputDataAnalysisToAgg) (int, error) {
	if cp == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp is nil")
	}
	wio, werr := gs3wfs(ctx, cp)
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
	wio.InputDataAnalysisToAgg = dataIn
	//md, err := json.Marshal(dataIn)
	//if err != nil {
	//	log.Err(err).Interface("dataIn", dataIn).Interface("wr", wr).Msg("SaveTaskOutput: failed")
	//	return 0, err
	//}
	//wr.Metadata = md
	err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	if dataIn.TextInput != nil {
		wio.PromptTextFromTextStage += *dataIn.TextInput
	}

	_, eerr := s3ws(ctx, cp, wio)
	if eerr != nil {
		log.Err(err).Msg("SaveTaskOutput: failed")
		return -1, eerr
	}
	return wr.WorkflowResultID, nil
}

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
	wio, werr := gs3wfs(ctx, cp)
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

package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/activity"
)

func WfDebugUtil(ctx context.Context, mb *MbChildSubProcessParams) error {
	zerr := S3WfRunUploadDebug(ctx, mb.GetRunName(), mb)
	if zerr != nil {
		log.Err(zerr).Msg("CsvIterator: SelectAiWorkflowAnalysisResultsIds failed")
		return zerr
	}
	return nil
}

func (z *ZeusAiPlatformActivities) CsvIterator(ctx context.Context, mb *MbChildSubProcessParams) error {
	werr := WfDebugUtil(ctx, mb)
	if werr != nil {
		return nil
	}
	in, gerr := gs3wfs(ctx, mb)
	if gerr != nil {
		log.Err(gerr).Msg("CsvIterator: gws failed")
		return gerr
	}
	log.Info().Interface("mb.Tc.TaskName", mb.Tc.TaskName).Msg("mb.Tc.TaskName)")
	log.Info().Interface("mb.Tc.TaskOffset", mb.Tc.TaskOffset).Msg("mb.Tc.TaskOffset)")
	sv, serr := artemis_orchestrations.SelectAiWorkflowAnalysisResultsIds(ctx, mb.Window, []int{mb.Oj.OrchestrationID}, []int{mb.Tc.TaskID}, mb.Tc.TaskOffset)
	if serr != nil {
		log.Err(serr).Msg("CsvIterator: SelectAiWorkflowAnalysisResultsIds failed")
		return serr
	}
	log.Info().Interface("len(sv)", len(sv)).Msg("CsvIterator")
	sm := make(map[int]map[int]bool)
	for _, vi := range sv {
		if _, ok := sm[vi.ChunkOffset]; !ok {
			sm[vi.ChunkOffset] = make(map[int]bool)
		}
		sm[vi.ChunkOffset][vi.IterationCount] = true
	}
	log.Info().Interface("sm", sm).Interface(" mb.Tc.ChunkIterator", mb.Tc.ChunkIterator).Interface("mb.Tc.TaskOffset", mb.Tc.TaskOffset).Msg("CsvIterator")
	prov := getPrompts(mb)
	log.Info().Interface("prov", prov).Msg("CsvIterator")
	for i := 0; i < mb.Tc.ChunkIterator; i++ {
		log.Info().Interface("i", i).Interface(" mb.Tc.ChunkIterator", mb.Tc.ChunkIterator).Interface("mb.Tc.TaskOffset", mb.Tc.TaskOffset).Msg("CsvIterator")
		err := iterResp(ctx, i, mb, in, prov, sm)
		if err != nil {
			log.Err(err).Msg("CsvIterator: gws failed")
			return err
		}
		activity.RecordHeartbeat(ctx, fmt.Sprintf("iterate-%d", i))
	}
	return nil
}

func iterResp(ctx context.Context, chunk int, mb *MbChildSubProcessParams, in *WorkflowStageIO, prms map[string]string, seen map[int]map[int]bool) error {
	// needs to get correct prompt mapped search
	sr := getSearchResults(chunk, mb, in)
	log.Info().Interface("len(sr)", len(sr)).Interface(" mb.Tc.ChunkIterator", mb.Tc.ChunkIterator).Msg("CsvIterator")
	var keys []string
	for key := range prms {
		keys = append(keys, key)
	}
	log.Info().Interface("keys", keys).Msg("iterResp")
	sort.Strings(keys)
	count := 0
	for offsetInd, colName := range keys {
		log.Info().Int("count", count).Msg("CsvIterator")
		count += 1
		if tv, ok := seen[chunk][offsetInd]; ok && tv {
			continue
		}
		log.Info().Interface("offsetInd", offsetInd).Interface("colName", colName).Msg("CsvIterator")
		taskInstPrompt := prms[colName]
		for _, v := range sr {
			if len(v.PromptKey) > 0 && v.PromptKey != colName {
				log.Warn().Interface("v.PromptKey", v.PromptKey).Interface("colName", colName).Msg("CsvIterator")
				continue
			}
			if !validPromptContent(ctx, v.Value) {
				log.Warn().Interface("v.Value", v.Value).Msg("CsvIterator")
				continue
			}
			na := NewZeusAiPlatformActivities()
			cr, err := na.CsvAnalysisTask(ctx, mb.Ou, getTaskPrompt(mb, taskInstPrompt), v.Value, true)
			if err != nil {
				log.Err(err).Msg("CsvIterator: iterResp failed")
				return err
			}
			log.Info().Str(fmt.Sprintf("iterate-%s-%d", colName, offsetInd), fmt.Sprintf("iterate-%s-%d", colName, offsetInd)).Msg("CsvIterator")
			activity.RecordHeartbeat(ctx, fmt.Sprintf("iterate-%s-%d", colName, offsetInd))
			err = saveCsvResp(ctx, colName, chunk, offsetInd, mb, cr, v)
			if err != nil {
				log.Err(err).Msg("CsvIterator: saveCsvResp failed")
				return err
			}
		}
	}
	return nil
}

func saveCsvResp(ctx context.Context, colName string, chunk, offsetInd int, mb *MbChildSubProcessParams, cr *ChatCompletionQueryResponse, v hera_search.SearchResult) error {
	log.Info().Msg("saveCsvResp")
	m := getCsvResp(colName, cr)
	if m == nil {
		log.Warn().Msg("saveCsvResp nil m")
		return nil
	}
	if len(v.QueryParams) > 0 {
		m["entity"] = strings.Join(v.QueryParams, ",")
	} else {
		m["entity"] = v.Source
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformActivities: saveCsvResp RecordCompletionResponse: failed")
		return err
	}
	wr := getWrAndIter(mb, chunk, offsetInd)
	rid, err := hera_openai_dbmodels.InsertCompletionResponseChatGpt(ctx, mb.Ou, cr.Response, b)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformActivities: saveCsvResp RecordCompletionResponse: failed")
	}
	wr.ResponseID = rid
	wr.TaskOffset = mb.Tc.TaskOffset
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("saveCsvResp: failed")
		return err
	}
	resp := InputDataAnalysisToAgg{
		CsvResponse: m,
	}
	err = s3wsCustomTaskName(ctx, mb, fmt.Sprintf("%d", wr.WorkflowResultID), resp)
	if err != nil {
		log.Err(err).Msg("s3wsCustomTaskName: saveCsvResp failed")
		return err
	}
	return nil
}

func getCsvResp(colName string, cr *ChatCompletionQueryResponse) map[string]interface{} {
	if cr == nil {
		return nil
	}
	if v, ok := cr.Prompt["response"]; ok {
		m := map[string]interface{}{
			colName: strings.TrimSpace(v),
		}
		return m
	}
	return nil
}

func validPromptContent(ctx context.Context, in string) bool {
	if len(in) <= 0 {
		return false
	}
	return true
}

const KevinFlowsOrgID = 1710298581127603000

func (z *ZeusAiPlatformActivities) CsvAnalysisTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, content string, isFlow bool) (*ChatCompletionQueryResponse, error) {
	cr := openai.ChatCompletionRequest{
		Model:       taskInst.AnalysisModel,
		Temperature: float32(taskInst.AnalysisTemperature),
		Messages:    []openai.ChatCompletionMessage{},
	}
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
	if taskInst.AnalysisMaxTokensPerTask > 0 {
		cr.MaxTokens = taskInst.AnalysisMaxTokensPerTask
	}
	prompt := make(map[string]string)
	prompt["prompt"] = taskInst.AnalysisPrompt
	prompt["content"] = content
	var oa hera_openai.OpenAI
	if isFlow {
		ou.OrgID = KevinFlowsOrgID
		ou.UserID = KevinFlowsOrgID
	}
	ps, err := GetMockingBirdSecrets(ctx, ou)
	if err != nil || ps == nil || ps.ApiKey == "" {
		if err == nil {
			err = fmt.Errorf("failed to get mockingbird secrets")
		}
		//log.Warn().Err(err).Msg("AiAnalysisTask: GetMockingbirdPlatformSecrets: failed to get mockingbird secrets")
		err = nil
		oa = hera_openai.HeraOpenAI
	} else {
		oa = hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
	}
	resp, err := oa.CreateChatCompletion(
		ctx, cr,
	)
	if err != nil {
		log.Err(err).Msg("CsvAnalysisTask")
		return nil, err
	}
	var sc string
	for _, c := range resp.Choices {
		sc += c.Message.Content + "\n"
	}
	prompt["response"] = sc
	return &ChatCompletionQueryResponse{
		Prompt:         prompt,
		ResponseTaskID: taskInst.AnalysisTaskID,
		Response:       resp,
	}, nil
}

func getPrompts(mb *MbChildSubProcessParams) map[string]string {
	if len(mb.WfExecParams.WorkflowOverrides.TaskPromptOverrides) == 0 {
		return nil
	}
	tmp := mb.WfExecParams.WorkflowOverrides.TaskPromptOverrides[mb.Tc.TaskName]
	return tmp.ReplacePrompts
}

func getSearchResults(chunk int, mb *MbChildSubProcessParams, in *WorkflowStageIO) []hera_search.SearchResult {
	sg := getJsonSgChunkToProcess2(chunk, mb, in)
	var sr []hera_search.SearchResult
	if sg.RegexSearchResults != nil {
		sr = sg.RegexSearchResults
	} else if sg.ApiResponseResults != nil {
		sr = sg.ApiResponseResults
	} else if in.PromptReduction.PromptReductionSearchResults != nil && in.PromptReduction.PromptReductionSearchResults.InSearchGroup != nil {
		sr = in.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults
	}
	return sr
}

func getTaskPrompt(mb *MbChildSubProcessParams, tp string) artemis_orchestrations.WorkflowTemplateData {
	fmt.Println(mb.Tc.TaskName)
	fmt.Println(mb.Tc.Model)
	fmt.Println(mb.Tc.ChunkIterator)
	taskInst := artemis_orchestrations.WorkflowTemplateData{
		AnalysisTaskDB: artemis_orchestrations.AnalysisTaskDB{
			AnalysisTaskID:         mb.Tc.TaskID,
			AnalysisModel:          mb.Tc.Model,
			AnalysisPrompt:         tp,
			AnalysisTemperature:    float64(mb.Tc.Temperature),
			AnalysisMarginBuffer:   mb.Tc.MarginBuffer,
			AnalysisTaskName:       mb.Tc.TaskName,
			AnalysisResponseFormat: mb.Tc.TaskType,
		},
	}
	return taskInst
}

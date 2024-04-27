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

func (z *ZeusAiPlatformActivities) CsvIterator(ctx context.Context, mb *MbChildSubProcessParams) error {
	in, gerr := gs3wfs(ctx, mb)
	if gerr != nil {
		log.Err(gerr).Msg("CsvIterator: gws failed")
		return gerr
	}
	log.Info().Interface("mb.Tc.TaskOffset", mb.Tc.TaskOffset).Msg("mb.Tc.TaskOffset)")
	sv, serr := artemis_orchestrations.SelectAiWorkflowAnalysisResultsIds(ctx, mb.Window, []int{mb.Oj.OrchestrationID}, []int{mb.Tc.TaskID}, mb.Tc.TaskOffset)
	if serr != nil {
		log.Err(serr).Msg("CsvIterator: SelectAiWorkflowAnalysisResultsIds failed")
		return serr
	}
	sm := make(map[int]map[int]bool)
	for _, vi := range sv {
		if _, ok := sm[vi.ChunkOffset]; !ok {
			sm[vi.ChunkOffset] = make(map[int]bool)
		}
		sm[vi.ChunkOffset][vi.IterationCount] = true
	}
	prov := getPrompts(mb)
	for i := 0; i < mb.Tc.ChunkIterator; i++ {
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
	sr := getSearchResults(chunk, mb, in)
	var keys []string
	for key := range prms {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for offsetInd, colName := range keys {
		if tv, ok := seen[chunk][offsetInd]; ok && tv {
			continue
		}
		taskInstPrompt := prms[colName]
		for _, v := range sr {
			if !validPromptContent(ctx, v.Value) {
				continue
			}
			na := NewZeusAiPlatformActivities()
			cr, err := na.CsvAnalysisTask(ctx, mb.Ou, getTaskPrompt(mb, taskInstPrompt), v.Value)
			if err != nil {
				log.Err(err).Msg("CsvIterator: iterResp failed")
				return err
			}
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

func (z *ZeusAiPlatformActivities) CsvAnalysisTask(ctx context.Context, ou org_users.OrgUser, taskInst artemis_orchestrations.WorkflowTemplateData, content string) (*ChatCompletionQueryResponse, error) {
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
	tmp := mb.WfExecParams.WorkflowOverrides.TaskPromptOverrides
	for _, mp := range tmp {
		return mp.ReplacePrompts
	}
	return nil
}

func getSearchResults(chunk int, mb *MbChildSubProcessParams, in *WorkflowStageIO) []hera_search.SearchResult {
	sg := getJsonSgChunkToProcess2(chunk, mb, in)
	var sr []hera_search.SearchResult
	if sg.RegexSearchResults != nil {
		sr = sg.RegexSearchResults
	} else if sg.ApiResponseResults != nil {
		sr = sg.ApiResponseResults
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

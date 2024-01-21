package ai_platform_service_orchestrations

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type FilteredMessages struct {
	MsgKeepIds []int `json:"msg_ids"`
}

func (z *ZeusAiPlatformServiceWorkflows) SocialMediaExtractionWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, sg *hera_openai_dbmodels.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		return nil, nil
	}
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, "ZeusAiPlatformServiceWorkflows", "SocialMediaExtractionWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	var aiResp *ChatCompletionQueryResponse
	extractCtx := workflow.WithActivityOptions(ctx, ao)
	switch sg.PlatformName {
	case twitterPlatform:
		err = workflow.ExecuteActivity(extractCtx, z.ExtractTweets, ou, sg).Get(extractCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run twitter extraction", "Error", err)
			return nil, err
		}
	case telegramPlatform:
	case discordPlatform:
	case redditPlatform:
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return aiResp, nil
}

func UnmarshallFilteredMsgIdsFromAiJson(fn string, cr *ChatCompletionQueryResponse) error {
	m, err := UnmarshallOpenAiJsonInterface(fn, cr)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
		return err
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: json.Marshal failed")
		return err
	}
	// Unmarshal the JSON string into the FilteredMessages struct
	cr.FilteredMessages = &FilteredMessages{}
	err = json.Unmarshal(jsonData, &cr.FilteredMessages)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: json.Unmarshal failed")
		return err
	}
	return nil
}

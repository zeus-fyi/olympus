package ai_platform_service_orchestrations

import (
	"fmt"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) SocialMediaMessagingWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, sg *hera_openai_dbmodels.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
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
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, "ZeusAiPlatformServiceWorkflows", "SocialMediaMessagingWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	var aiResp *ChatCompletionQueryResponse
	extractCtx := workflow.WithActivityOptions(ctx, ao)
	saveAiResponseCtx := workflow.WithActivityOptions(ctx, ao)

	switch sg.PlatformName {
	case twitterPlatform:
		err = workflow.ExecuteActivity(extractCtx, z.CreateTweets, ou, sg).Get(extractCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run twitter extraction", "Error", err)
			return nil, err
		}
		if aiResp == nil {
			return nil, fmt.Errorf("no ai response")
		}
		err = workflow.ExecuteActivity(saveAiResponseCtx, z.RecordCompletionResponse, ou, aiResp).Get(saveAiResponseCtx, nil)
		if err != nil {
			logger.Error("failed to save completion response", "Error", err)
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

func UnmarshallTwitterFromAiJsonSlice(fn string, cr *ChatCompletionQueryResponse) ([]*twitter.CreateTweetRequest, error) {
	m, err := UnmarshallOpenAiJsonInterface(fn, cr)
	if err != nil {
		return nil, err
	}
	rb, ok := m[text].([]string)
	if !ok || len(rb) == 0 {
		return nil, fmt.Errorf("text body had no text, or was not a string")
	}

	rt, ok := m[inReplyToTweetID].([]string)
	var crts []*twitter.CreateTweetRequest
	tr := &twitter.CreateTweetRequest{}
	for i, v := range rb {
		tr = &twitter.CreateTweetRequest{
			Text: v,
		}
		if ok && len(rt) >= i {
			tr.Reply = &twitter.CreateTweetReply{InReplyToTweetID: rt[i]}
		}
		if ok {

		}
		crts = append(crts, tr)
	}

	return crts, nil
}

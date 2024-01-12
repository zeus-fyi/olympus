package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) SendResponseToApiForScoresInJson(ctx context.Context, cr *ChatCompletionQueryResponse) (map[string]interface{}, error) {
	return nil, nil
}

func (z *ZeusAiPlatformActivities) EvalLookup(ctx context.Context, ou org_users.OrgUser, evalID int) ([]artemis_orchestrations.EvalFn, error) {
	evalFn, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(ctx, ou, evalID)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return nil, err
	}
	return evalFn, nil
}

func (z *ZeusAiPlatformActivities) SaveEvalMetricResults(ctx context.Context, emr *artemis_orchestrations.EvalMetricsResults) error {
	if emr == nil || emr.EvalMetricsResults == nil {
		log.Info().Msg("SaveEvalMetricResults: emr is nil")
		return nil
	}
	err := artemis_orchestrations.UpsertEvalMetricsResults(ctx, emr.EvalContext, emr.EvalMetricsResults)
	if err != nil {
		log.Err(err).Msg("EvalLookup: failed to get eval fn")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) CreateJsonOutputModelResponse(ctx context.Context, ou org_users.OrgUser, params hera_openai.OpenAIParams) (*ChatCompletionQueryResponse, error) {
	var err error
	var resp openai.ChatCompletionResponse
	ps, err := GetMockingBirdSecrets(ctx, ou)
	if err != nil || ps == nil || ps.ApiKey == "" {
		log.Info().Msg("CreatJsonOutputModelResponse: GetMockingBirdSecrets failed to find user openai api key, using system key")
		err = nil
		resp, err = hera_openai.HeraOpenAI.MakeCodeGenRequestJsonFormattedOutput(ctx, ou, params)
	} else {
		oc := hera_openai.InitOrgHeraOpenAI(ps.ApiKey)
		resp, err = oc.MakeCodeGenRequestJsonFormattedOutput(ctx, ou, params)
	}
	if err != nil {
		log.Err(err).Msg("CreatJsonOutputModelResponse: MakeCodeGenRequestJsonFormattedOutput failed")
		return nil, err
	}
	return &ChatCompletionQueryResponse{
		Prompt:   map[string]string{"prompt": params.Prompt},
		Response: resp,
	}, nil
}

func (z *ZeusAiPlatformActivities) EvalModelScoredJsonOutput(ctx context.Context, jsonInterface map[string]interface{}, evalFn *artemis_orchestrations.EvalFn) (*artemis_orchestrations.EvalMetricsResults, error) {
	if len(jsonInterface) == 0 || evalFn == nil {
		log.Info().Msg("EvalModelScoredJsonOutput: at least one input is nil or empty")
		return nil, nil
	}
	scoredResults, err := TransformJSONToEvalScoredMetrics(jsonInterface, evalFn.EvalMetricMap)
	if err != nil {
		log.Err(err).Msg("EvalModelScoredJsonOutput: failed to transform json to eval scored metrics")
		return nil, err
	}
	return scoredResults, nil
}

// TransformEvalMetricsToJSONSchema transforms a slice of EvalMetric into jsonschema.Definition.
func TransformEvalMetricsToJSONSchema(metrics []artemis_orchestrations.EvalMetric) (jsonschema.Definition, error) {
	properties := make(map[string]jsonschema.Definition)
	var required []string
	for _, metric := range metrics {
		dataType, itemsDef, err := mapDataType(metric.EvalMetricDataType)
		if err != nil {
			return jsonschema.Definition{}, err
		}
		propertyDef := jsonschema.Definition{
			Type:        dataType,
			Description: metric.EvalModelPrompt,
		}
		if itemsDef != nil {
			propertyDef.Items = itemsDef
		}
		properties[metric.EvalMetricName] = propertyDef
		required = append(required, metric.EvalMetricName)
	}

	return jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   required,
	}, nil
}

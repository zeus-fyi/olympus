package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) AnalyzeEngagementTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		log.Info().Msg("AnalyzeEngagementTweets: no search results to analyze engagement")
		return nil, nil
	}
	za := NewZeusAiPlatformActivities()
	params := hera_openai.OpenAIParams{
		Model:              sg.Model,
		FunctionDefinition: sg.FunctionDefinition,
		Prompt:             hera_search.FormatSearchResultsV3(sg.SearchResults),
	}
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, params)
	if err != nil {
		log.Err(err).Msg("AnalyzeEngagementTweets: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	m, err := UnmarshallOpenAiJsonInterface(sg.FunctionDefinition.Name, resp)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
		return nil, err
	}
	jsd := artemis_orchestrations.ConvertToJsonSchema(sg.FunctionDefinition)
	if err != nil {
		log.Err(err).Msg("AnalyzeEngagementTweets: mapToSchema failed")
		return nil, err
	}
	resp.JsonResponseResults = jsd
	return resp, nil
}

func FilterAndExtractRelevantTweetsJsonSchemaFunctionDef3() openai.FunctionDefinition {
	fdSchema := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"analyzed_tweet": {
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"msg_id": {
						Type:        jsonschema.Number,
						Description: "The system message ID of the analyzed tweet which is given by the user json key msg_id",
					},
					"lead_score": {
						Type:        jsonschema.Integer,
						Description: scoreInst,
					},
				},
				Required: []string{"analyzed_msg_id", "lead_score"},
			},
		},
		Required: []string{"analyzed_tweet"},
	}
	fd := openai.FunctionDefinition{
		Name:       "twitter_engagement_scores",
		Parameters: fdSchema,
	}
	return fd
}

const scoreInst = `Create a lead score using the below metrics Scoring Metrics:
Error Mention (429, 5xx, etc.): 10 points. High importance as it directly relates to the problem the
load balancer solves. Mention of Dissatisfaction with Current RPC Solution: 8 points.
Indicates a clear need for a new solution. Engagement or mentions related relevant content: 5 points.
Suggests interest in the topic. Organization Size/Market Potential: Score based on potential usage volume
(e.g., large enterprises 10 points, mid-size businesses 5 points, small businesses 3 points).
If you are unable to tell, set it to 0. Mention of EVM/Blockchain Transaction Challenges: 10 points. 
Indicating direct relevance and need. Interest in EVM tx simulations or Smart Contracts: 8 points. Shows alignment 
with the product offering. Developer Platform Discussion: Assign 10 points for any mention of "developer platform" or similar terms
that align with Zeusfyi's target audience. Kubernetes Engagement: Allot 15 points when Kubernetes is mentioned, indicating
a direct interest or existing investment in the technology that Zeusfyi enhances. Complex Deployment Solutions: Add 10 points for
discussions about challenging deployments, indicating a need for Zeusfyi's deployment strategies. Cloud Management and Flexibility: 
Score 8 points for expressions of a need for versatile cloud management, which is central to Zeusfyi’s capabilities. Multi-Cloud and 
Multi-Tenancy: Give 12 points for any reference to multi-cloud strategies or multi-tenancy challenges, as Zeusfyi excels in these areas. 
Temporal Workflow Interest: Allocate 7 points for interests in temporal workflows, suggesting a good fit for Zeusfyi's process orchestration tools.
AI-Powered Optimization: Add 10 points for mentions of AI or iterative solver workflows, as this is a unique feature of Zeusfyi that can offer significant value to the user.`

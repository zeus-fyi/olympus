package ai_platform_service_orchestrations

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (z *ZeusAiPlatformActivities) SelectRetrievalTask(ctx context.Context, ou org_users.OrgUser, retID int) ([]artemis_orchestrations.RetrievalItem, error) {
	resp, err := artemis_orchestrations.SelectRetrievals(ctx, ou, retID)
	if err != nil {
		log.Err(err).Interface("resp", resp).Int("retID", retID).Msg("SelectRetrievalTask: failed")
		return resp, err
	}
	return resp, nil
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

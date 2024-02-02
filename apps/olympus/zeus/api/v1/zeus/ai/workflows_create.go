package zeus_v1_ai

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

type PostWorkflowsRequest struct {
	WorkflowName          string                                   `json:"workflowName"`
	WorkflowGroupName     string                                   `json:"workflowGroupName"`
	StepSize              int                                      `json:"stepSize"`
	StepSizeUnit          string                                   `json:"stepSizeUnit"`
	Models                TaskMap                                  `json:"models"`
	EvalsMap              map[string]artemis_orchestrations.EvalFn `json:"evalsMap,omitempty"`
	EvalTasksMap          TaskEvalsMap                             `json:"evalTasksMap,omitempty"`
	AggregateSubTasksMap  AggregateSubTasksMap                     `json:"aggregateSubTasksMap,omitempty"`
	AnalysisRetrievalsMap AnalysisRetrievalsMap                    `json:"analysisRetrievalsMap"`
}

type TaskEvalsMap map[string]map[string]bool
type AnalysisRetrievalsMap map[string]map[string]bool
type AggregateSubTasksMap map[string]map[string]bool
type TaskMap map[string]TaskModelInstructions

// TaskModelInstructions represents the equivalent of the TypeScript interface TaskModelInstructions
type TaskModelInstructions struct {
	TaskID                int                             `json:"taskID"`
	Model                 string                          `json:"model"`
	TaskType              string                          `json:"taskType"`
	TaskGroup             string                          `json:"taskGroup"`
	TaskName              string                          `json:"taskName"`
	MaxTokens             int                             `json:"maxTokens"`
	TokenOverflowStrategy string                          `json:"tokenOverflowStrategy"`
	Prompt                string                          `json:"prompt"`
	CycleCount            int                             `json:"cycleCount"`
	EvalFns               []artemis_orchestrations.EvalFn `json:"evalFns,omitempty"`
}

func PostWorkflowsRequestHandler(c echo.Context) error {
	request := new(PostWorkflowsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateWorkflow(c)
}

func (w *PostWorkflowsRequest) CreateOrUpdateWorkflow(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if w.WorkflowName == "" || len(w.Models) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if err != nil {
		log.Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}
	if w.WorkflowName == "" || len(w.Models) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	wt := artemis_orchestrations.WorkflowTemplate{
		WorkflowName:              w.WorkflowName,
		WorkflowGroup:             w.WorkflowGroupName,
		FundamentalPeriod:         w.StepSize,
		FundamentalPeriodTimeUnit: w.StepSizeUnit,
	}
	wft := artemis_orchestrations.WorkflowTasks{
		AggTasks:          []artemis_orchestrations.AggTask{},
		AnalysisOnlyTasks: []artemis_orchestrations.AITaskLibrary{},
	}

	ms := make(map[string]string)
	for _, m := range w.Models {
		if m.CycleCount < 1 {
			m.CycleCount = 1
		}
		switch m.TaskType {
		case "aggregation":
			agt := artemis_orchestrations.AggTask{
				AggId:      m.TaskID,
				CycleCount: m.CycleCount,
				Tasks:      []artemis_orchestrations.AITaskLibrary{},
			}
			if w.EvalTasksMap != nil {
				if evm, tok := w.EvalTasksMap[fmt.Sprintf("%d", m.TaskID)]; tok {
					for k, v := range evm {
						if v {
							mappedEval := w.EvalsMap[fmt.Sprintf("%d", k)]
							if mappedEval.EvalStrID != nil && *mappedEval.EvalStrID != "" {
								eid, serr := strconv.Atoi(*mappedEval.EvalStrID)
								if serr != nil {
									log.Err(serr).Msg("failed to parse int")
									return c.JSON(http.StatusBadRequest, nil)
								}
								mappedEval.EvalID = aws.Int(eid)
							}
							agt.EvalFns = append(agt.EvalFns, mappedEval)
						}
					}
				}
			}
			for k, v := range w.AggregateSubTasksMap {
				if k == fmt.Sprintf("%d", m.TaskID) {
					for at, isTrue := range v {
						if isTrue {
							ms[at] = fmt.Sprintf("%d", m.TaskID)
							ait, zerr := strconv.Atoi(at)
							if zerr != nil {
								log.Err(zerr).Msg("failed to parse int")
								return c.JSON(http.StatusBadRequest, nil)
							}
							ta := artemis_orchestrations.AITaskLibrary{
								TaskStrID:  at,
								TaskID:     ait,
								CycleCount: m.CycleCount,
							}
							if w.EvalTasksMap != nil {
								if evm, tok := w.EvalTasksMap[at]; tok {
									for ke, ve := range evm {
										if ve {
											mappedEval := w.EvalsMap[ke]
											if mappedEval.EvalStrID != nil && *mappedEval.EvalStrID != "" {
												eid, serr := strconv.Atoi(*mappedEval.EvalStrID)
												if serr != nil {
													log.Err(serr).Msg("failed to parse int")
													return c.JSON(http.StatusBadRequest, nil)
												}
												mappedEval.EvalID = aws.Int(eid)
											}
											agt.EvalFns = append(agt.EvalFns, mappedEval)
										}
									}
								}
							}
							agt.Tasks = append(agt.Tasks, ta)
						}
					}
				}
			}
			wft.AggTasks = append(wft.AggTasks, agt)
		case "analysis":
			at := artemis_orchestrations.AITaskLibrary{
				TaskID:                m.TaskID,
				OrgID:                 ou.OrgID,
				UserID:                ou.UserID,
				MaxTokensPerTask:      m.MaxTokens,
				TaskType:              m.TaskType,
				TaskName:              m.TaskName,
				TaskGroup:             m.TaskGroup,
				TokenOverflowStrategy: m.TokenOverflowStrategy,
				Model:                 m.Model,
				Prompt:                m.Prompt,
				CycleCount:            m.CycleCount,
				RetrievalDependencies: []artemis_orchestrations.RetrievalItem{},
			}

			if w.EvalTasksMap != nil {
				if evm, tok := w.EvalTasksMap[fmt.Sprintf("%d", m.TaskID)]; tok {
					for ke, ve := range evm {
						if ve {
							mappedEval := w.EvalsMap[ke]
							if mappedEval.EvalStrID != nil && *mappedEval.EvalStrID != "" {
								eid, serr := strconv.Atoi(*mappedEval.EvalStrID)
								if serr != nil {
									log.Err(serr).Msg("failed to parse int")
									return c.JSON(http.StatusBadRequest, nil)
								}
								mappedEval.EvalID = aws.Int(eid)
							}
							at.EvalFns = append(at.EvalFns, mappedEval)
						}
					}
				}
			}
			for _, v := range w.AnalysisRetrievalsMap {
				for rt, isTrue := range v {
					if isTrue {
						rid, rerr := strconv.Atoi(rt)
						if rerr != nil {
							log.Err(rerr).Msg("failed to parse int")
							return c.JSON(http.StatusBadRequest, nil)
						}
						at.RetrievalDependencies = append(at.RetrievalDependencies, artemis_orchestrations.RetrievalItem{
							RetrievalID: &rid,
						})
					}
				}
			}
			wft.AnalysisOnlyTasks = append(wft.AnalysisOnlyTasks, at)
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	err = artemis_orchestrations.InsertWorkflowWithComponents(c.Request().Context(), ou, &wt, wft)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, wt)
}

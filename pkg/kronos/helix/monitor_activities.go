package kronos_helix

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	IrisHealthEndpoint      = "https://iris.zeus.fyi/health"
	HestiaHealthEndpoint    = "https://hestia.zeus.fyi/health"
	ZeusHealthEndpoint      = "https://api.zeus.fyi/health"
	ZeusCloudHealthEndpoint = "https://cloud.zeus.fyi"
)

type MonitorInstructions struct {
	ServiceName           string        `json:"serviceName"`
	Endpoint              string        `json:"endpoint"`
	PollInterval          time.Duration `json:"pollInterval"`
	AlertFailureThreshold int           `json:"alertThreshold"`
}

func CreateNewMonitorInstructions(serviceName, endpoint string, pollInterval time.Duration, alertThreshold int) MonitorInstructions {
	return MonitorInstructions{
		ServiceName:           serviceName,
		Endpoint:              endpoint,
		PollInterval:          pollInterval,
		AlertFailureThreshold: alertThreshold,
	}
}

func CreateHealthMonitorAlertEvent(serviceName string) *pagerduty.V2Event {
	return &pagerduty.V2Event{
		Action: apollo_pagerduty.TRIGGER,
		Payload: &pagerduty.V2Payload{
			Summary:   fmt.Sprintf("There is an unhealthy service: %s", serviceName),
			Source:    "HEALTH_ALERTS",
			Severity:  apollo_pagerduty.CRITICAL,
			Component: fmt.Sprintf("This is a microservice %s component", serviceName),
			Group:     "This is a microservice group",
			Class:     "Microservice Health Monitoring",
			Details:   nil,
		},
	}
}
func CalculatePollCycles(intervalResetTime time.Duration, pollInterval time.Duration) int {
	return int(intervalResetTime / pollInterval)
}

func (k *KronosWorkflow) Monitor(ctx workflow.Context, oj *artemis_orchestrations.OrchestrationJob, mi Instructions, pollCycles int) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    mi.Monitors.PollInterval,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	failureCount := 0
	for i := 0; i < pollCycles; i++ {
		healthCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(healthCtx, k.CheckEndpointHealth, mi, failureCount).Get(healthCtx, &failureCount)
		if err != nil {
			logger.Error("failed to execute triggered alert", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
		if failureCount >= mi.Monitors.AlertFailureThreshold {
			alertCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(alertCtx, k.ExecuteTriggeredAlert, CreateHealthMonitorAlertEvent(mi.Monitors.ServiceName)).Get(alertCtx, nil)
			if err != nil {
				logger.Error("failed to execute triggered alert", "Error", err)
				// You can decide if you want to return the error or continue monitoring.
				return err
			}
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				TaskQueue:         KronosHelixTaskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}
			childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
			childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "OrchestrationChildProcessReset", &oj, mi)
			var childWE workflow.Execution
			if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
				logger.Error("Failed to get child workflow execution", "Error", err)
				return err
			}
			return nil
		}
		err = workflow.Sleep(ctx, mi.Monitors.PollInterval)
		if err != nil {
			logger.Error("failed to sleep", "Error", err)
			return err
		}
	}
	return nil
}

func (k *KronosActivities) CheckEndpointHealth(ctx context.Context, mi Instructions, failures int) (int, error) {
	resp, err := http.Get(mi.Monitors.Endpoint)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("CheckEndpointHealth: get health check failed")
		return failures + 1, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return failures + 1, nil
	}
	return 0, nil
}

package kronos_helix

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/rs/zerolog/log"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
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

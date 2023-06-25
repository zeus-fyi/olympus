package hestia_ovhcloud

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

var priceMapHourly = map[string]float64{
	"b2-7":   0.0813,
	"b2-15":  0.1539,
	"b2-30":  0.3123,
	"b2-60":  0.606,
	"b2-120": 1.1923,
	"c2-7":   0.1176,
	"c2-15":  0.2276,
	"c2-30":  0.4586,
	"c2-60":  0.8986,
	"c2-120": 1.7786,
	"r2-15":  0.1176,
	"r2-30":  0.1363,
	"r2-60":  0.2639,
	"r2-120": 0.5323,
	"r2-240": 1.05,
}

var priceMapMonthly = map[string]float64{
	"b2-7":   29.04,
	"b2-15":  55.44,
	"b2-30":  112.20,
	"b2-60":  217.80,
	"b2-120": 429.00,
	"c2-7":   42.24,
	"c2-15":  81.83,
	"c2-30":  165.00,
	"c2-60":  323.40,
	"c2-120": 640.20,
	"r2-15":  42.24,
	"r2-30":  48.84,
	"r2-60":  95.03,
	"r2-120": 191.40,
	"r2-240": 376.20,
}

type OvhInstanceType struct {
	Category string `json:"category"`
	Gpus     int    `json:"gpus"`
	Name     string `json:"name"`
	Ram      int    `json:"ram"`
	State    string `json:"state"`
	VCPUs    int    `json:"vCPUs"`

	Details Details `json:"details"`
}

type Details struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	Ram               int    `json:"ram"`
	Disk              int    `json:"disk"`
	Vcpus             int    `json:"vcpus"`
	Type              string `json:"type"`
	OsType            string `json:"osType"`
	InboundBandwidth  int    `json:"inboundBandwidth"`
	OutboundBandwidth int    `json:"outboundBandwidth"`
	Available         bool   `json:"available"`
	PlanCodes         struct {
		Monthly string `json:"monthly"`
		Hourly  string `json:"hourly"`
	} `json:"planCodes"`
	Capabilities []struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	} `json:"capabilities"`
	Quota int `json:"quota"`
}

func (o *OvhCloud) GetSizes(ctx context.Context, serviceName, region string) ([]OvhInstanceType, error) {
	endpoint := fmt.Sprintf("/cloud/project/%s/capabilities/kube/flavors?region=%s", serviceName, region)
	var instances []OvhInstanceType
	err := o.Get(endpoint, &instances)
	if err != nil {
		log.Err(err).Msg("OvhCloud: GetSizes")
		return nil, err
	}
	return instances, nil
}

func (o *OvhCloud) GetFlavorDetails(ctx context.Context, serviceName, region string) ([]Details, error) {
	endpoint := fmt.Sprintf("/cloud/project/%s/flavor?region=%s", serviceName, region)
	var instances []Details
	err := o.GetWithContext(ctx, endpoint, &instances)
	if err != nil {
		log.Err(err).Msg("OvhCloud: GetSizes")
		return nil, err
	}
	return instances, nil
}

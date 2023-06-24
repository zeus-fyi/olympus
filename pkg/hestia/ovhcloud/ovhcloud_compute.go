package hestia_ovhcloud

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type OvhInstanceType struct {
	Category string `json:"category"`
	Gpus     int    `json:"gpus"`
	Name     string `json:"name"`
	Ram      int    `json:"ram"`
	State    string `json:"state"`
	VCPUs    int    `json:"vCPUs"`
}

// curl -X GET "/v1/cloud/project/serviceName/capabilities/kube/flavors?region=US-WEST-OR-1" \
// func (c *Client) CallAPIWithContext(ctx context.Context, method string, path string, reqBody interface{}, resType interface{}, needAuth bool) error

func (o *OvhCloud) GetSizes(ctx context.Context, serviceName, region string) ([]OvhInstanceType, error) {
	endpoint := fmt.Sprintf("/v1/cloud/project/%s/capabilities/kube/flavors?region=%s", serviceName, region)
	var instances []OvhInstanceType
	err := o.CallAPIWithContext(ctx, "GET", endpoint, "", instances, true)
	if err != nil {
		log.Err(err).Msg("OvhCloud: GetSizes")
		return nil, err
	}
	return instances, nil
}

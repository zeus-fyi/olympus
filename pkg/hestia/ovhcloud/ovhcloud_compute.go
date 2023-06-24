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

package hestia_digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_ovhcloud "github.com/zeus-fyi/olympus/pkg/hestia/ovhcloud"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"k8s.io/apimachinery/pkg/api/errors"
)

// https://docs.digitalocean.com/reference/api/api-reference/#tag/Domain-Records

const (
	Sfo3LoadBalancerIp = "143.198.244.181"
	NycLoadBalancerIp  = "164.90.252.115"
	GkeUsCentral1Ip    = "34.122.201.76"

	OvhUsWestOr1InternalLoadBalancerIp = "51.81.201.5"
	OvhUsWestOr1ExternalLoadBalancerIp = "51.81.201.16"

	AwsUsWest1Ip = "a48dca93c6961441e8829cc9e99fd21a-a52026aa1784e018.elb.us-west-1.amazonaws.com."
)

func (d *DigitalOcean) CreateDomain(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) (*godo.DomainRecord, error) {
	loadBalancer := ""
	createRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: cloudCtxNs.Namespace,
		Data: loadBalancer,
		TTL:  3600,
	}
	switch cloudCtxNs.CloudProvider {
	case "ovh":
		switch cloudCtxNs.Context {
		case hestia_ovhcloud.OvhInternalContext:
			loadBalancer = OvhUsWestOr1InternalLoadBalancerIp
		case hestia_ovhcloud.OvhSharedContext:
			loadBalancer = OvhUsWestOr1ExternalLoadBalancerIp
		}
	case "gcp":
		loadBalancer = GkeUsCentral1Ip
	case "do":
		switch cloudCtxNs.Region {
		case "nyc1":
			loadBalancer = NycLoadBalancerIp
		case "sfo3":
			loadBalancer = Sfo3LoadBalancerIp
		}
	case "aws":
		createRequest.Type = "CNAME"
		createRequest.TTL = 43200
		loadBalancer = AwsUsWest1Ip
	}
	createRequest.Data = loadBalancer
	dr, _, err := d.Domains.CreateRecord(ctx, "zeus.fyi", createRequest)
	if errors.IsAlreadyExists(err) {
		return nil, err
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create domain record")
		return dr, err
	}
	return dr, err
}

func (d *DigitalOcean) RemoveSubDomainARecord(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	recordType := "A"
	if cloudCtxNs.CloudProvider == "aws" {
		recordType = "CNAME"
	}
	dl, _, err := d.Domains.RecordsByType(ctx, "zeus.fyi", recordType, &godo.ListOptions{})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create domain record")
		return err
	}
	for _, dn := range dl {
		if dn.Name == cloudCtxNs.Namespace {
			fmt.Println("deleting", dn.Name, dn.ID)
			_, err = d.Domains.DeleteRecord(ctx, "zeus.fyi", dn.ID)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("failed to create domain record")
				return err
			}
		}
	}
	return err
}

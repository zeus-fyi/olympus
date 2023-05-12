package hestia_digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

// https://docs.digitalocean.com/reference/api/api-reference/#tag/Domain-Records

const (
	Sfo3LoadBalancerIp = "143.198.244.181"
	NycLoadBalancerIp  = "164.90.252.115"
	GkeUsCentral1Ip    = "34.122.201.76"

	AwsUsWest1Ip = "a48dca93c6961441e8829cc9e99fd21a-a52026aa1784e018.elb.us-west-1.amazonaws.com"
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
		loadBalancer = AwsUsWest1Ip
	}
	createRequest.Data = loadBalancer

	dr, _, err := d.Domains.CreateRecord(ctx, "zeus.fyi", createRequest)
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

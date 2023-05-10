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
)

func (d *DigitalOcean) CreateDomain(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) (*godo.DomainRecord, error) {
	loadBalancer := ""
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
	}
	createRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: cloudCtxNs.Namespace,
		Data: loadBalancer,
		TTL:  3600,
	}

	dr, _, err := d.Domains.CreateRecord(ctx, "zeus.fyi", createRequest)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create domain record")
		return dr, err
	}
	return dr, err
}

func (d *DigitalOcean) RemoveSubDomainARecord(ctx context.Context, name string) error {
	dl, _, err := d.Domains.RecordsByType(ctx, "zeus.fyi", "A", &godo.ListOptions{})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create domain record")
		return err
	}
	for _, dn := range dl {
		if dn.Name == name {
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

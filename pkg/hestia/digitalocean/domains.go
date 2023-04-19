package hestia_digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

// https://docs.digitalocean.com/reference/api/api-reference/#tag/Domain-Records

const NycLoadBalancerIp = "164.90.252.115"

func (d *DigitalOcean) CreateDomain(ctx context.Context, name string) (*godo.DomainRecord, error) {
	createRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: name,
		Data: NycLoadBalancerIp,
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
	sbName := fmt.Sprintf("%s.zeus.fyi", name)
	for _, dn := range dl {
		if dn.Name == sbName {
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

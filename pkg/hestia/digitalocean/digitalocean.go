package hestia_digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

type DigitalOcean struct {
	*godo.Client
}

func InitDoClient(ctx context.Context, token string) DigitalOcean {
	do := godo.NewFromToken(token)
	return DigitalOcean{do}
}

func (d *DigitalOcean) GetSizes(ctx context.Context) ([]godo.Size, error) {
	//lo := godo.ListOptions{
	//	Page:         0,
	//	PerPage:      0,
	//	WithProjects: false,
	//}
	sizes, _, err := d.Sizes.List(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
	}
	return sizes, err
}

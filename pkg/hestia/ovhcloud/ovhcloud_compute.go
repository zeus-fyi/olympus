package hestia_ovhcloud

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (o *OvhCloud) GetSizes(ctx context.Context) (any, error) {
	// func (c *Client) CallAPIWithContext(ctx context.Context, method string, path string, reqBody interface{}, resType interface{}, needAuth bool) error
	err := o.CallAPIWithContext(ctx, "", "", "", "", true)
	if err != nil {
		log.Err(err).Msg("OvhCloud: GetSizes")
		return nil, err
	}
	return nil, nil
}

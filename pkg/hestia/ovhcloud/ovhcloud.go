package hestia_ovhcloud

import (
	"context"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
)

const (
	OvhUS = "ovh-us"
)

type OvhCloud struct {
	*ovh.Client
}

type OvhCloudCreds struct {
	Region      string
	AppKey      string
	AppSecret   string
	ConsumerKey string
}

func InitOvhClient(ctx context.Context, creds OvhCloudCreds) OvhCloud {
	if creds.Region == "" {
		creds.Region = OvhUS
	}
	client, err := ovh.NewClient(
		creds.Region,
		creds.AppKey,
		creds.AppSecret,
		creds.ConsumerKey,
	)
	if err != nil {
		panic(err)
	}
	return OvhCloud{client}
}

func (o *OvhCloud) GetSizes(ctx context.Context) error {
	// func (c *Client) CallAPIWithContext(ctx context.Context, method string, path string, reqBody interface{}, resType interface{}, needAuth bool) error
	err := o.CallAPIWithContext(ctx, "", "", "", "", true)
	if err != nil {
		log.Err(err).Msg("OvhCloud: GetSizes")
		return err
	}
	return nil
}

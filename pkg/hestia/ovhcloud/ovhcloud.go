package hestia_ovhcloud

import (
	"context"

	"github.com/ovh/go-ovh/ovh"
)

const (
	OvhUS                  = "ovh-us"
	OvhRegionUsWestOr1     = "us-west-or-1"
	OvhRegionUsWestOr1ENUM = "US-WEST-OR-1"
	OvhServiceName         = "e39851e46915473bb1c50dd56f987c26"
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

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
	OvhInternalKubeID      = "750cf38b-0965-4b2b-b6ba-9728ca3f239e"
	OvhSharedKubeID        = "a7ea8ded-fa8f-48f3-83d7-ce01410552bc"

	OvhSharedContext   = "zeusfyi-shared"
	OvhInternalContext = "zeusfyi"
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

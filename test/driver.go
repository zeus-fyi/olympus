package test

import (
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

var host = "http://localhost:9001/v1/infra/create"

var kns = autok8s_core.KubeCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	CtxType:       "dev-sfo3-zeus",
	Namespace:     "demo",
}

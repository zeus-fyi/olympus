package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
)

type Spec struct {
	DefaultBackend   *common.ParentClass
	IngressClassName *common.ParentClass
	TLS              TLS
	Rules            common.ParentClass
}

func NewIngressSpec() Spec {
	spec := Spec{
		DefaultBackend:   nil,
		IngressClassName: nil,
		TLS:              NewTLS(),
		Rules:            common.ParentClass{},
	}
	return spec
}

package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
)

type Spec struct {
	DefaultBackend   *common.ParentClass
	IngressClassName *common.ParentClass
	TLS              common.ParentClass
	Rules            common.ParentClass
}

func NewIngressSpec() Spec {
	// TODO give these default parent names
	spec := Spec{
		DefaultBackend:   nil,
		IngressClassName: nil,
		TLS:              common.ParentClass{},
		Rules:            common.ParentClass{},
	}
	return spec
}

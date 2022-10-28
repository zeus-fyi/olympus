package ingresses

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

type Spec struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	DefaultBackend   *structs.SuperParentClassGroup
	IngressClassName *IngressClassName
	TLS              TLS
	Rules            Rules
}

func NewIngressSpec() Spec {
	spec := Spec{
		DefaultBackend:   nil,
		IngressClassName: nil,
		TLS:              NewTLS(),
		Rules:            NewRules(),
	}
	spec.ChartSubcomponentParentClassTypeName = "Spec"
	spec.ChartComponentResourceID = IngressChartComponentResourceID
	return spec
}

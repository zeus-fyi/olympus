package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

type ServicePort struct {
	Values structs.ChildValuesMap
}
type ServicePorts []ServicePort

func NewServicePorts() ServicePorts {
	sps := []ServicePort{NewServicePort()}
	return sps
}

func NewServicePort() ServicePort {
	sp := ServicePort{}
	fields := []string{"name", "protocol", "port", "targetPort", "nodePort"}
	s := structs.NewChildValuesMapKeyFromIterable(fields...)
	sp.Values = s
	return sp
}

/* Kubernetes Reference
type ServicePort struct {
	Name       string
	Protocol   string
	Port       int
	TargetPort string
	NodePort   int
}
*/

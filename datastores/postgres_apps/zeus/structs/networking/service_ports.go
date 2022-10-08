package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
)

type ServicePort struct {
	Values common.ChildValuesMap
}
type ServicePorts []ServicePort

func NewServicePorts() ServicePorts {
	sps := []ServicePort{NewServicePort()}
	return sps
}

func NewServicePort() ServicePort {
	sp := ServicePort{}
	fields := []string{"name", "protocol", "port", "targetPort", "nodePort"}
	s := common.NewChildValuesMapKeyFromIterable(fields...)
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

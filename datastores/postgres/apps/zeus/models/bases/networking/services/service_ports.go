package services

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

type ServicePorts struct {
	Ports []structs.ChildClassMultiValue
}

type Port structs.ChildClassMultiValue

func NewServicePorts() ServicePorts {
	sp := ServicePorts{Ports: []structs.ChildClassMultiValue{}}
	return sp
}

func (ss *ServiceSpec) AddPortMapValuesThenInsertAsPort(svcPortMap map[string]string) {
	if len(ss.Ports) <= 0 {
		ss.Ports = []structs.ChildClassMultiValue{}
	}
	port := structs.ChildClassMultiValue{}
	port.ChartSubcomponentParentClassTypeID = ss.ChartSubcomponentParentClassTypeID
	port.ChartSubcomponentChildClassTypeName = "svcPort"
	for k, v := range svcPortMap {
		if k == "name" {
			port.ChartSubcomponentChildClassTypeName += k
		}
		port.Values = append(port.Values, autogen_bases.ChartSubcomponentsChildValues{ChartSubcomponentKeyName: k, ChartSubcomponentValue: v})
	}
	ss.Ports = append(ss.Ports, port)
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

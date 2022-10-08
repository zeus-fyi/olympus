package networking

import (
	v1 "k8s.io/api/core/v1"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

func ServicePortsToDB(cps []v1.ServicePort) networking.ServicePorts {
	spSlice := make(networking.ServicePorts, len(cps))
	for i, p := range cps {
		port := ServicePortToDB(p)
		spSlice[i] = port
	}
	return spSlice
}

func ServicePortToDB(p v1.ServicePort) networking.ServicePort {
	sp := networking.ServicePort{
		Name:       p.Name,
		Protocol:   string(p.Protocol),
		Port:       int(p.Port),
		TargetPort: p.TargetPort.String(),
		NodePort:   int(p.NodePort),
	}
	return sp
}

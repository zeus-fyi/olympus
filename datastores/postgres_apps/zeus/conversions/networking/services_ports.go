package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/common"
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

func ServicePortToDB(svcPort v1.ServicePort) networking.ServicePort {
	sp := networking.NewServicePort()
	sp.Values["name"] = common.ConvertKeyValueStringToChildValues("name", svcPort.Name)
	sp.Values["port"] = common.ConvertKeyValueStringToChildValues("port", string(svcPort.Port))
	sp.Values["targetPort"] = common.ConvertKeyValueStringToChildValues("targetPort", svcPort.TargetPort.String())
	sp.Values["nodePort"] = common.ConvertKeyValueStringToChildValues("nodePort", string(svcPort.NodePort))
	return sp
}

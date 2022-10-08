package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/networking"
	v1 "k8s.io/api/core/v1"
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

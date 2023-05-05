package services

import (
	"fmt"
)

func (s *Service) ServicePortsToDB() {
	for _, svcPort := range s.K8sService.Spec.Ports {
		m := make(map[string]string)
		m["port"] = fmt.Sprintf("%d", svcPort.Port)
		if len(svcPort.Name) > 0 {
			m["name"] = svcPort.Name
		}
		tgtPortStr := svcPort.TargetPort.String()
		if len(tgtPortStr) > 0 && tgtPortStr != "0" {
			m["targetPort"] = svcPort.TargetPort.String()
		}
		if svcPort.NodePort != 0 {
			m["nodePort"] = fmt.Sprintf("%d", svcPort.NodePort)
		}
		if len(string(svcPort.Protocol)) > 0 {
			m["protocol"] = string(svcPort.Protocol)
		}
		s.ServiceSpec.AddPortMapValuesThenInsertAsPort(m)
	}
}

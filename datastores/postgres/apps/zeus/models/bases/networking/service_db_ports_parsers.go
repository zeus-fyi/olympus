package networking

import "fmt"

func (s *Service) ServicePortsToDB() {
	for _, svcPort := range s.K8sService.Spec.Ports {
		m := make(map[string]string)
		m["name"] = svcPort.Name
		m["port"] = fmt.Sprintf("%d", svcPort.Port)
		m["targetPort"] = svcPort.TargetPort.String()
		m["nodePort"] = fmt.Sprintf("%d", svcPort.NodePort)
		m["protocol"] = string(svcPort.Protocol)
		s.ServiceSpec.AddPortMapValuesThenInsertAsPort(m)
	}
}

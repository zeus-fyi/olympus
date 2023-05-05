package services

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *Service) ParseSvcPorts(portMap map[int][]common_conversions.PC) error {
	for _, port := range portMap {
		k8sPort := v1.ServicePort{}
		for _, portComponent := range port {
			val := portComponent.ChartSubcomponentValue
			switch portComponent.ChartSubcomponentKeyName {
			case "targetPort":
				if val != "" && val != "0" {
					k8sPort.TargetPort = intstr.FromString(val)
				}
			case "port":
				k8sPort.Port = string_utils.ConvertStringTo32BitInt(val)
			case "name":
				k8sPort.Name = val
			case "nodePort":
				k8sPort.NodePort = string_utils.ConvertStringTo32BitInt(val)
			case "protocol":
				k8sPort.Protocol = v1.Protocol(val)
			}
		}
		s.K8sService.Spec.Ports = append(s.K8sService.Spec.Ports, k8sPort)
	}
	return nil
}

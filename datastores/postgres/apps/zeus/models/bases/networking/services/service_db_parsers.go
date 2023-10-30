package services

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
)

func (s *Service) ConvertK8sServiceToDB() {
	s.Metadata.ChartSubcomponentParentClassTypeName = "ServiceParentMetadata"
	metadata := common_conversions.ConvertMetadata(s.K8sService.ObjectMeta)
	s.Metadata.Metadata = metadata
	s.ConvertServiceSpecConfigToDB()
}

func (s *Service) ConvertServiceSpecConfigToDB() {
	s.ServiceSpec = NewServiceSpec()
	s.ServiceSpec.Selector = common_conversions.ConvertSelectorByFields(s.K8sService.Spec.Selector)
	s.Type.ChartSubcomponentValue = string(s.K8sService.Spec.Type)
	s.ClusterIP.ChartSubcomponentValue = s.K8sService.Spec.ClusterIP
	s.ServicePortsToDB()
	return
}

package servicemonitors

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func (s *ServiceMonitor) ConvertK8sServiceMonitorToDB() error {
	s.Metadata.ChartSubcomponentParentClassTypeName = "ServiceMonitorParentMetadata"
	s.Metadata.Metadata = common_conversions.ConvertMetadata(s.K8sServiceMonitor.ObjectMeta)
	s.Metadata.ChartComponentResourceID = ServiceMonitorChartComponentResourceID
	err := s.ConvertK8sServiceMonitorSpecToDB()
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

func (s *ServiceMonitor) ConvertK8sServiceMonitorSpecToDB() error {
	bytes, err := json.Marshal(s.K8sServiceMonitor.Spec)
	if err != nil {
		log.Err(err)
		return err
	}
	s.Spec.ChildClassSingleValue = structs.NewInitChildClassSingleValue("spec", string(bytes))
	s.Spec.ChartSubcomponentParentClassTypeName = "Spec"
	s.Spec.ChildClassSingleValue.ChartSubcomponentChildClassTypeName = "ServiceMonitorSpec"
	return nil
}

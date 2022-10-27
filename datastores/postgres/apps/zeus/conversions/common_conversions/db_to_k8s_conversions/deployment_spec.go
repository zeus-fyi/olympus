package db_to_k8s_conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentSpec(k8sDeploymentSpec *v1.DeploymentSpec, pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "replicas":
			k8sDeploymentSpec.Replicas = string_utils.ConvertStringTo32BitPtrInt(pc.ChartSubcomponentValue)
		case "selectorString":
			err := ParseLabelSelectorJsonString(k8sDeploymentSpec.Selector, pc.ChartSubcomponentValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

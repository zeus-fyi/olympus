package configuration

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (cm *ConfigMap) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Data":
			// TODO
			err := cm.ParseCmDataToK8ss(pc)
			if err != nil {
				return err
			}
		case "ConfigMapParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&cm.K8sConfigMap.ObjectMeta, pc)
		}
	}
	return nil
}

func (cm *ConfigMap) ParseCmDataToK8ss(pcSlice []common_conversions.PC) error {
	// TODO
	return nil
}

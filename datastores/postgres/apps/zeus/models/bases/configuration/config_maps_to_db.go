package configuration

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"

func (cm *ConfigMap) ConvertK8sConfigMapToDB() {
	metadata := common_conversions.ConvertMetadata(cm.K8sConfigMap.ObjectMeta)
	cm.Metadata.Metadata = metadata
	cm.Data = NewCMData(cm.K8sConfigMap.Data)
}

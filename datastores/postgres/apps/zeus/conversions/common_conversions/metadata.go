package common_conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertMetadata(m metav1.ObjectMeta) structs.Metadata {
	dbMetaConfig := structs.NewMetadata()
	dbMetaConfig.Name.ChartSubcomponentValue = m.Name
	dbMetaConfig.Annotations.Values = ConvertKeyValueToChildValues(m.Annotations)
	dbMetaConfig.Labels.Values = ConvertKeyValueToChildValues(m.Labels)
	return dbMetaConfig
}

func CreateMetadataByFields(name string, annotations, labels map[string]string) structs.Metadata {
	dbMetaConfig := structs.NewMetadata()
	dbMetaConfig.Name.ChartSubcomponentValue = name
	dbMetaConfig.Annotations.Values = ConvertKeyValueToChildValues(annotations)
	dbMetaConfig.Labels.Values = ConvertKeyValueToChildValues(labels)
	return dbMetaConfig
}

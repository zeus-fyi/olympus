package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertMetadata(m metav1.ObjectMeta) common.Metadata {
	dbMetaConfig := common.NewMetadata()
	dbMetaConfig.Name.ChartSubcomponentValue = m.Name
	dbMetaConfig.Annotations.AnnotationValues = ConvertKeyValueToChildValues(m.Annotations)
	dbMetaConfig.Labels.LabelValues = ConvertKeyValueToChildValues(m.Labels)
	return dbMetaConfig
}

func CreateMetadataByFields(name string, annotations, labels map[string]string) common.Metadata {
	dbMetaConfig := common.NewMetadata()
	dbMetaConfig.Name.AddNameValue(name)
	dbMetaConfig.Annotations.AnnotationValues = ConvertKeyValueToChildValues(annotations)
	dbMetaConfig.Labels.LabelValues = ConvertKeyValueToChildValues(labels)
	return dbMetaConfig
}

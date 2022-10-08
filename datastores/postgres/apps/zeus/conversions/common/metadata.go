package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertMetadata(m metav1.ObjectMeta) common.Metadata {
	dbMetaConfig := common.NewMetadata()
	dbMetaConfig.Name = m.Name
	dbMetaConfig.Annotations = ConvertKeyValueToChildValues(m.Annotations)
	dbMetaConfig.Labels = ConvertKeyValueToChildValues(m.Labels)
	return dbMetaConfig
}

func CreateMetadataByFields(name string, annotations, labels map[string]string) common.Metadata {
	dbMetaConfig := common.NewMetadata()
	dbMetaConfig.Name = name
	dbMetaConfig.Annotations = ConvertKeyValueToChildValues(annotations)
	dbMetaConfig.Labels = ConvertKeyValueToChildValues(labels)
	return dbMetaConfig
}

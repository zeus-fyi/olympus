package common_conversions

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertMetadata(m metav1.ObjectMeta) structs.Metadata {
	dbMetaConfig := structs.NewMetadata()
	dbMetaConfig.Name.ChartSubcomponentValue = m.Name
	dbMetaConfig.Annotations.Values = ConvertKeyValueToChildValues(m.Annotations)
	m.Labels = AddVersionIDLabel(m.Labels)
	dbMetaConfig.Labels.Values = ConvertKeyValueToChildValues(m.Labels)
	return dbMetaConfig
}

func CreateMetadataByFields(name string, annotations, labels map[string]string) structs.Metadata {
	dbMetaConfig := structs.NewMetadata()
	dbMetaConfig.Name.ChartSubcomponentValue = name
	dbMetaConfig.Annotations.Values = ConvertKeyValueToChildValues(annotations)
	AddVersionIDLabel(labels)
	dbMetaConfig.Labels.Values = ConvertKeyValueToChildValues(labels)
	return dbMetaConfig
}

func AddVersionIDLabel(labels map[string]string) map[string]string {
	if len(labels) <= 0 {
		labels = make(map[string]string)
	}
	var ts chronos.Chronos
	labels["version"] = fmt.Sprintf("version-%d", ts.UnixTimeStampNow())
	return labels
}

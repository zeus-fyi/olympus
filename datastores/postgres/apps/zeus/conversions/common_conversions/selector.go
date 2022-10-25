package common_conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertSelector(m *metav1.LabelSelector) structs.Selector {
	dbSelectorConfig := structs.NewSelector()
	dbSelectorConfig.MatchLabels.Values = ConvertKeyValueToChildValues(m.MatchLabels)
	return dbSelectorConfig
}

func ConvertSelectorByFields(labels map[string]string) structs.Selector {
	dbSelectorConfig := structs.NewSelector()
	dbSelectorConfig.MatchLabels.Values = ConvertKeyValueToChildValues(labels)
	return dbSelectorConfig
}

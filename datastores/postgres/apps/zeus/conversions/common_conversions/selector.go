package common_conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertSelector(m *metav1.LabelSelector) common.Selector {
	dbSelectorConfig := common.NewSelector()
	dbSelectorConfig.MatchLabels.Values = ConvertKeyValueToChildValues(m.MatchLabels)
	return dbSelectorConfig
}

func ConvertSelectorByFields(labels map[string]string) common.Selector {
	dbSelectorConfig := common.NewSelector()
	dbSelectorConfig.MatchLabels.Values = ConvertKeyValueToChildValues(labels)
	return dbSelectorConfig
}

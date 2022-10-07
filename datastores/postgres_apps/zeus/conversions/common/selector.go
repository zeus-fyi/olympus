package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertSelector(m *metav1.LabelSelector) common.Selector {
	dbSelectorConfig := common.NewSelector()
	dbSelectorConfig.MatchLabels = ConvertKeyValueToChildValues(m.MatchLabels)
	return dbSelectorConfig
}

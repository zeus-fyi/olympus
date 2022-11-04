package chart_workload

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (nk *NativeK8s) IdWorkloadFromBytes(jsonBytes []byte) (metav1.TypeMeta, error) {
	metaType := metav1.TypeMeta{}

	m := make(map[string]interface{})
	err := json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return metaType, err
	}
	for k, v := range m {
		switch k {
		case "kind":
			metaType.Kind = v.(string)
		case "apiVersion":
			metaType.APIVersion = v.(string)
		}
	}
	return metaType, err
}

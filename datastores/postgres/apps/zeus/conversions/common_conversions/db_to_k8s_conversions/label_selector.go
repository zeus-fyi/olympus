package db_to_k8s_conversions

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseLabelSelectorJsonString(selectorString string) (*metav1.LabelSelector, error) {
	selectorLabel := metav1.LabelSelector{}

	var m map[string]interface{}
	err := json.Unmarshal([]byte(selectorString), &m)
	if err != nil {
		return &selectorLabel, err
	}

	bytes, berr := getBytes(m)
	if berr != nil {
		return &selectorLabel, berr
	}
	perr := json.Unmarshal(bytes, &selectorLabel)
	if perr != nil {
		return &selectorLabel, perr
	}

	return &selectorLabel, nil
}

func getBytes(structToBytes interface{}) ([]byte, error) {
	bytes, berr := json.Marshal(structToBytes)
	return bytes, berr
}

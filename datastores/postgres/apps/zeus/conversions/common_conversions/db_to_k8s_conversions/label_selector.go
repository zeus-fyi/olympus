package db_to_k8s_conversions

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseLabelSelectorJsonString(selectorLabel *metav1.LabelSelector, selectorString string) error {
	bytes, berr := getBytes(selectorString)
	if berr != nil {
		return berr
	}
	perr := json.Unmarshal(bytes, &selectorLabel)
	if perr != nil {
		return perr
	}

	return nil
}

func getBytes(structToBytes interface{}) ([]byte, error) {
	bytes, berr := json.Marshal(structToBytes)
	return bytes, berr
}

package db_to_k8s_conversions

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseLabelSelectorJsonString(selectorString string) (*metav1.LabelSelector, error) {
	selectorLabel := metav1.LabelSelector{}

	var m map[string]interface{}
	err := json.Unmarshal([]byte(selectorString), &m)
	if err != nil {
		log.Err(err).Msg("ParseLabelSelectorJsonString Unmarshal selectorString")
		return &selectorLabel, err
	}

	bytes, berr := getBytes(m)
	if berr != nil {
		log.Err(err).Msg("ParseLabelSelectorJsonString getBytes")
		return &selectorLabel, berr
	}
	perr := json.Unmarshal(bytes, &selectorLabel)
	if perr != nil {
		log.Err(perr).Msg("ParseLabelSelectorJsonString Unmarshal")
		return &selectorLabel, perr
	}

	return &selectorLabel, nil
}

func getBytes(structToBytes interface{}) ([]byte, error) {
	bytes, berr := json.Marshal(structToBytes)
	return bytes, berr
}

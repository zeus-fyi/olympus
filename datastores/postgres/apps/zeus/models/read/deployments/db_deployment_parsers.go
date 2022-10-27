package read_deployments

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseDeploymentParentChildAggValues(ckaggString string) error {
	m := make(map[string][]map[string]interface{})
	err := json.Unmarshal([]byte(ckaggString), &m)
	if err != nil {
		return err
	}
	for k, _ := range m {
		if k == "parentWrapper" {
			switch k {
			case "DeploymentParentMetadata":
				//_, _ = ParseMetadataValues(bytesN)
			}
		}
	}
	return nil
}

func ParseMetadataValues(metadataStringBytes []byte) (metav1.ObjectMeta, error) {
	metaData := metav1.ObjectMeta{}
	err := json.Unmarshal(metadataStringBytes, &metaData)
	if err != nil {
		return metaData, err
	}
	return metaData, nil
}

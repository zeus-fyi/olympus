package conversions

import (
	"encoding/json"

	v1 "k8s.io/api/apps/v1"

	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

var yr = transformations.YamlReader{}

func SaveDeploymentConfigToDB() error {
	filepath := "/Users/alex/Desktop/Zeus/olympus/pkg/zeus/core/transformations/deployment.yaml"
	jsonBytes, err := yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	return err
}

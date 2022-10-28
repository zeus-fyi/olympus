package create_deployments

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) (deployments.Deployment, error) {
	dbDeployment := deployments.NewDeployment()
	dbDeployment.Metadata.Metadata = common_conversions.CreateMetadataByFields(d.Name, d.Annotations, d.Labels)
	depSpec, err := ConvertDeploymentSpec(d.Spec)
	if err != nil {
		return dbDeployment, err
	}
	dbDeployment.Spec = depSpec
	return dbDeployment, nil
}

func ConvertDeploymentSpec(ds v1.DeploymentSpec) (deployments.Spec, error) {
	deploymentTemplateSpec := ds.Template
	podTemplateSpec := deploymentTemplateSpec.Spec

	dbDeploymentSpec := deployments.NewDeploymentSpec()

	m := make(map[string]string)
	if ds.Selector != nil {
		bytes, err := json.Marshal(ds.Selector)
		if err != nil {
			return dbDeploymentSpec, err
		}
		selectorString := string(bytes)
		m["selectorString"] = selectorString
		dbDeploymentSpec.Selector.MatchLabels.AddValues(m)
	}

	dbDeploymentSpec.Replicas.ChartSubcomponentValue = string_utils.Convert32BitPtrIntToString(ds.Replicas)
	dbPodTemplateSpec, err := dbDeploymentSpec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return dbDeploymentSpec, err
	}
	dbDeploymentSpec.Template = dbPodTemplateSpec
	return dbDeploymentSpec, nil
}

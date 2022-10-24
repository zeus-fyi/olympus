package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) (Deployment, error) {
	dbDeployment := NewDeployment()
	dbDeployment.Metadata.Metadata = common_conversions.CreateMetadataByFields(d.Name, d.Annotations, d.Labels)
	depSpec, err := ConvertDeploymentSpec(d.Spec)
	if err != nil {
		return dbDeployment, err
	}
	dbDeployment.Spec.DeploymentSpec = depSpec
	return dbDeployment, nil
}

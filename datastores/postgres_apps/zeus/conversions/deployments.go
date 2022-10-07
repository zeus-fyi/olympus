package conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/workloads"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) workloads.Deployment {
	dbDeployment := workloads.NewDeployment()
	dbDeployment.Metadata = CreateMetadataByFields(d.Name, d.Annotations, d.Labels)
	dbDeployment.Spec = ConvertDeploymentSpec(d.Spec)
	return dbDeployment
}

func ConvertDeploymentSpec(ds v1.DeploymentSpec) workloads.DeploymentSpec {
	deploymentTemplateSpec := ds.Template
	podTemplateSpec := deploymentTemplateSpec.Spec
	dbPodTemplateSpec := ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	dbDeploymentSpec := workloads.DeploymentSpec{
		Replicas: 0,
		Template: dbPodTemplateSpec,
	}
	return dbDeploymentSpec
}

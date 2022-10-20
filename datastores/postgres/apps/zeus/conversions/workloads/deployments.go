package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) workloads.Deployment {
	dbDeployment := workloads.NewDeployment()
	dbDeployment.Metadata.Metadata = common.CreateMetadataByFields(d.Name, d.Annotations, d.Labels)
	dbDeployment.Spec.DeploymentSpec = ConvertDeploymentSpec(d.Spec)
	return dbDeployment
}

func ConvertDeploymentSpec(ds v1.DeploymentSpec) workloads.DeploymentSpec {
	deploymentTemplateSpec := ds.Template
	podTemplateSpec := deploymentTemplateSpec.Spec
	dbPodTemplateSpec := containers.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	dbDeploymentSpec := workloads.DeploymentSpec{
		// TODO Replicas: ,
		Template: dbPodTemplateSpec,
		Selector: common.ConvertSelector(ds.Selector),
	}
	return dbDeploymentSpec
}

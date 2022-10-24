package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) (workloads.Deployment, error) {
	dbDeployment := workloads.NewDeployment()
	dbDeployment.Metadata.Metadata = common.CreateMetadataByFields(d.Name, d.Annotations, d.Labels)
	depSpec, err := ConvertDeploymentSpec(d.Spec)
	if err != nil {
		return dbDeployment, err
	}
	dbDeployment.Spec.DeploymentSpec = depSpec
	return dbDeployment, nil
}

func ConvertDeploymentSpec(ds v1.DeploymentSpec) (workloads.DeploymentSpec, error) {
	deploymentTemplateSpec := ds.Template
	podTemplateSpec := deploymentTemplateSpec.Spec
	dbDeploymentSpec := workloads.DeploymentSpec{
		// TODO Replicas: ,
		Selector: common.ConvertSelector(ds.Selector),
	}
	dbPodTemplateSpec, err := containers.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return dbDeploymentSpec, err
	}
	dbDeploymentSpec.Template = dbPodTemplateSpec
	return dbDeploymentSpec, nil
}

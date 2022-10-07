package conversions

import (
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	v1 "k8s.io/api/apps/v1"
)

func ConvertDeploymentConfigToDB(d *v1.Deployment) error {
	err := ConvertDeploymentSpec(&d.Spec)
	return err
}

func ConvertDeploymentSpec(ds *v1.DeploymentSpec) error {
	deploymentTemplateSpec := ds.Template
	podTemplateSpec := deploymentTemplateSpec.Spec

	dbPodTemplateSpec := ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	err := dev_hacks.Use(dbPodTemplateSpec)
	return err
}

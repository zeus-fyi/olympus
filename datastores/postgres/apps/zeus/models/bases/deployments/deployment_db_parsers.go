package deployments

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (d *Deployment) ParsePCGroupMap(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := db_to_k8s_conversions.ConvertDeploymentSpec(&d.K8sDeployment.Spec, pc)
			if err != nil {
				return err
			}
		case "DeploymentParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&d.K8sDeployment.ObjectMeta, pc)
		}
	}
	return nil
}

func (d *Deployment) ConvertDeploymentConfigToDB() error {
	dbDeployment := NewDeployment()
	dbDeployment.Metadata.Metadata = common_conversions.CreateMetadataByFields(d.K8sDeployment.Name, d.K8sDeployment.Annotations, d.K8sDeployment.Labels)
	err := d.ConvertDeploymentSpec()
	if err != nil {
		return err
	}
	return nil
}

func (d *Deployment) ConvertDeploymentSpec() error {
	deploymentTemplateSpec := d.K8sDeployment.Spec.Template
	podTemplateSpec := deploymentTemplateSpec.Spec

	dbDeploymentSpec := NewDeploymentSpec()

	m := make(map[string]string)
	if d.K8sDeployment.Spec.Selector != nil {
		bytes, err := json.Marshal(d.K8sDeployment.Spec.Selector)
		if err != nil {
			return err
		}
		selectorString := string(bytes)
		m["selectorString"] = selectorString
		dbDeploymentSpec.Selector.MatchLabels.AddValues(m)
	}

	dbDeploymentSpec.Replicas.ChartSubcomponentValue = string_utils.Convert32BitPtrIntToString(d.K8sDeployment.Spec.Replicas)
	dbPodTemplateSpec, err := dbDeploymentSpec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return err
	}
	dbDeploymentSpec.Template = dbPodTemplateSpec
	return nil
}

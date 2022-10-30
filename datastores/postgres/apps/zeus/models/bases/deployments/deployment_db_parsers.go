package deployments

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (d *Deployment) ConvertDeploymentConfigToDB() error {
	newDeployment := NewDeployment()
	d.KindDefinition = newDeployment.KindDefinition
	d.Metadata = newDeployment.Metadata
	d.Spec = newDeployment.Spec
	d.Metadata.Metadata = common_conversions.CreateMetadataByFields(d.K8sDeployment.Name, d.K8sDeployment.Annotations, d.K8sDeployment.Labels)
	err := d.ConvertDeploymentSpec()
	if err != nil {
		return err
	}
	return nil
}

func (d *Deployment) ConvertDeploymentSpec() error {
	dbDeploymentSpec := NewDeploymentSpec()
	m := make(map[string]string)
	if d.K8sDeployment.Spec.Selector != nil {
		bytes, err := json.Marshal(d.K8sDeployment.Spec.Selector)
		if err != nil {
			return err
		}
		selectorString := string(bytes)
		m["selectorString"] = selectorString
		d.Spec.Selector.MatchLabels.AddValues(m)
	}

	d.Spec.Replicas.ChartSubcomponentValue = string_utils.Convert32BitPtrIntToString(d.K8sDeployment.Spec.Replicas)
	dbPodTemplateSpec, err := dbDeploymentSpec.Template.ConvertPodTemplateSpecConfigToDB(&d.K8sDeployment.Spec.Template.Spec)
	dbPodTemplateSpecMetadata := d.K8sDeployment.Spec.Template.GetObjectMeta()
	dbPodTemplateSpec.Metadata.Metadata = common_conversions.CreateMetadataByFields(dbPodTemplateSpecMetadata.GetName(), dbPodTemplateSpecMetadata.GetAnnotations(), dbPodTemplateSpecMetadata.GetLabels())

	if err != nil {
		return err
	}
	d.Spec.Template = dbPodTemplateSpec
	return nil
}

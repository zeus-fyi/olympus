package deployments

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (d *Deployment) ParsePCGroupMap(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := d.ConvertDBDeploymentSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "PodTemplateSpec":
			db_to_k8s_conversions.ConvertPodSpecField(&d.K8sDeployment.Spec.Template, pc)
		case "DeploymentParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&d.K8sDeployment.ObjectMeta, pc)
		case "PodTemplateSpecMetadata":
			db_to_k8s_conversions.ConvertMetadata(&d.K8sDeployment.Spec.Template.ObjectMeta, pc)
		}
	}
	return nil
}

func (d *Deployment) ConvertDBDeploymentSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "replicas":
			d.K8sDeployment.Spec.Replicas = string_utils.ConvertStringTo32BitPtrInt(pc.ChartSubcomponentValue)
		case "selector":
			sl, err := db_to_k8s_conversions.ParseLabelSelectorJsonString(pc.ChartSubcomponentValue)
			d.K8sDeployment.Spec.Selector = sl
			if err != nil {
				log.Err(err).Msg("ConvertDBDeploymentSpecToK8s")
				return err
			}
		}
	}
	return nil
}

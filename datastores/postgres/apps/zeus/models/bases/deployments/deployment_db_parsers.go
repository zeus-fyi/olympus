package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
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

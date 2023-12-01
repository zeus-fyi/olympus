package read_deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	v1 "k8s.io/api/core/v1"
)

func DBDeploymentResource(d *deployments.Deployment, ckagg, podSpecVolumesStr string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		return pcerr
	}
	pcerr = d.ParsePCGroupMap(pcGroupMap)
	if pcerr != nil {
		return pcerr
	}
	if len(podSpecVolumesStr) > 0 {
		vs, vserr := db_to_k8s_conversions.ParsePodSpecDBVolumesString(podSpecVolumesStr)
		if vserr != nil {
			return vserr
		}
		d.K8sDeployment.Spec.Template.Spec.Volumes = vs

	}
	return nil
}

func DBDeploymentContainer(d *deployments.Deployment, c *containers.Container) error {
	if c.Metadata.ContainerID != 0 {
		cerr := c.ParseFields()
		if cerr != nil {
			return cerr
		}
		deploymentContainers := d.K8sDeployment.Spec.Template.Spec.Containers
		if len(deploymentContainers) <= 0 {
			deploymentContainers = []v1.Container{}
		}
		c.K8sContainer.Name = c.Metadata.ContainerName
		c.K8sContainer.Image = c.Metadata.ContainerImageID
		c.K8sContainer.ImagePullPolicy = v1.PullPolicy(c.Metadata.ContainerImagePullPolicy)

		if c.Metadata.IsInitContainer {
			if d.K8sDeployment.Spec.Template.Spec.InitContainers == nil {
				d.K8sDeployment.Spec.Template.Spec.InitContainers = []v1.Container{}
			}
			d.K8sDeployment.Spec.Template.Spec.InitContainers = append(d.K8sDeployment.Spec.Template.Spec.InitContainers, c.K8sContainer)
		} else {
			d.K8sDeployment.Spec.Template.Spec.Containers = append(deploymentContainers, c.K8sContainer)
		}
	}
	return nil
}

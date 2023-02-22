package read_statefulsets

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
	v1 "k8s.io/api/core/v1"
)

func DBStatefulSetResource(sts *statefulsets.StatefulSet, ckagg, podSpecVolumesStr string) error {
	pcGroupMap, pcerr := common_conversions.ParseDeploymentParentChildAggValues(ckagg)
	if pcerr != nil {
		return pcerr
	}
	pcerr = sts.ParseDBConfigToK8s(pcGroupMap)
	if pcerr != nil {
		return pcerr
	}
	if len(podSpecVolumesStr) > 0 {
		vs, vserr := db_to_k8s_conversions.ParsePodSpecDBVolumesString(podSpecVolumesStr)
		if vserr != nil {
			return vserr
		}
		sts.K8sStatefulSet.Spec.Template.Spec.Volumes = vs

	}
	return nil
}

func DBStatefulSetContainer(d *statefulsets.StatefulSet, c *containers.Container) error {
	if c.Metadata.ContainerID != 0 {
		cerr := c.ParseFields()
		if cerr != nil {
			return cerr
		}
		c.K8sContainer.Name = c.Metadata.ContainerName
		c.K8sContainer.Image = c.Metadata.ContainerImageID
		c.K8sContainer.ImagePullPolicy = v1.PullPolicy(c.Metadata.ContainerImagePullPolicy)

		if c.Metadata.IsInitContainer {
			deploymentInitContainers := d.K8sStatefulSet.Spec.Template.Spec.InitContainers
			if len(deploymentInitContainers) <= 0 {
				deploymentInitContainers = []v1.Container{}
			}
			d.K8sStatefulSet.Spec.Template.Spec.InitContainers = append(deploymentInitContainers, c.K8sContainer)
		} else {
			deploymentContainers := d.K8sStatefulSet.Spec.Template.Spec.Containers
			if len(deploymentContainers) <= 0 {
				deploymentContainers = []v1.Container{}
			}
			d.K8sStatefulSet.Spec.Template.Spec.Containers = append(deploymentContainers, c.K8sContainer)
		}
	}
	return nil
}

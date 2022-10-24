package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

func (d *Deployment) GetContainers() containers.Containers {
	return d.Spec.Template.Spec.PodTemplateContainers
}

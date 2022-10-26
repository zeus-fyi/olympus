package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
)

func NewPodSpec() PodTemplateSpec {
	return PodTemplateSpec{}
}

type PodTemplateSpec struct {
	containers.PodTemplateSpec
	ContainerMap map[int]Container
}

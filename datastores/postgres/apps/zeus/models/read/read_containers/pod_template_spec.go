package read_containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	v1 "k8s.io/api/core/v1"
)

func NewPodSpec() PodTemplateSpec {
	return PodTemplateSpec{}
}

type PodTemplateSpec struct {
	containers.PodTemplateSpec
	ContainerMap map[int]Container
}

func ParsePodSpecDBVolumesString(podSpecVolumes string) ([]v1.Volume, error) {
	m := make(map[string]interface{})
	var vms []v1.Volume
	err := json.Unmarshal([]byte(podSpecVolumes), &m)
	if err != nil {
		return vms, err
	}
	for _, v := range m {
		bytes, berr := json.Marshal(v)
		if berr != nil {
			return vms, berr
		}
		var vm v1.Volume
		perr := json.Unmarshal(bytes, &vm)
		if perr != nil {
			return vms, perr
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

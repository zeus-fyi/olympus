package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) insertContainerVolumeMountsValues(m map[string]containers.Container, imageID string, contVmsSubCTE, contVmsRelationshipsSubCTE *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	ts := chronos.Chronos{}
	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()

	for _, vm := range c.VolumeMounts {
		vmID := ts.UnixTimeStampNow()
		contVmsSubCTE.AddValues(vm.VolumeReadOnly, vm.VolumeSubPath, vmID, vm.VolumeMountPath, vm.VolumeName)
		contVmsRelationshipsSubCTE.AddValues(podSpecChildClassTypeID, c.GetContainerID(), vmID)
	}
	return
}

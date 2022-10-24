package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodContainersGroup) insertContainerVolumeMountsValues(podSpecChildClassTypeID int, imageID string, contVmsSubCTE, contVmsRelationshipsSubCTE *sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	ts := chronos.Chronos{}
	for _, vm := range c.VolumeMounts {
		vmID := ts.UnixTimeStampNow()
		contVmsSubCTE.AddValues(vmID, vm.VolumeMountPath, vm.VolumeName)
		contVmsRelationshipsSubCTE.AddValues(podSpecChildClassTypeID, selectRelatedContainerIDFromImageID(imageID), vmID)
	}
	return
}

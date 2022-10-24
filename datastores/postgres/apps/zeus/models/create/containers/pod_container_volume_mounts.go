package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodContainersGroup) insertContainerVolumeMountsHeader() string {
	return "INSERT INTO container_volume_mounts(volume_mount_id, volume_mount_path, volume_name) VALUES "
}

func (p *PodContainersGroup) getInsertContainerVolumeMountsValues(imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for _, vm := range c.VolumeMounts {
		cteSubfield.AddValues(vm.VolumeMountID, vm.VolumeMountPath, vm.VolumeName)
	}
	return
}

func (p *PodContainersGroup) insertContainerVolumeMountRelationshipHeader() string {
	return "INSERT INTO containers_volume_mounts(chart_subcomponent_child_class_type_id, container_id, volume_mount_id) VALUES "
}

func (p *PodContainersGroup) getContainerVolumeMountRelationshipValues(podSpecChildClassTypeID int, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for _, vm := range c.VolumeMounts {
		cteSubfield.AddValues(podSpecChildClassTypeID, selectRelatedContainerIDFromImageID(imageID), vm.VolumeMountID)
	}
	return
}

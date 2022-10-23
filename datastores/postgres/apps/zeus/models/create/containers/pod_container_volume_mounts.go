package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

func (p *PodContainersGroup) insertContainerVolumeMountsHeader() string {
	return "INSERT INTO container_volume_mounts(probe_id, probe_key_values_jsonb) VALUES "
}

func (p *PodContainersGroup) getInsertContainerVolumeMountsValues(parentExpression string, contVolMounts autogen_bases.ContainerVolumeMountsSlice, isLastValuesGroup bool) string {
	for i, vm := range contVolMounts {
		parentExpression += fmt.Sprintf("\n('%d', '%s', '%s')", vm.VolumeMountID, vm.VolumeMountPath, vm.VolumeMountPath)
		if i < len(contVolMounts)-1 && !isLastValuesGroup {
			parentExpression += ","
		}
	}
	return parentExpression
}

func (p *PodContainersGroup) insertContainerVolumeMountRelationshipHeader() string {
	return "INSERT INTO containers_volume_mounts(chart_subcomponent_child_class_type_id, container_id, volume_mount_id) VALUES "
}

func (p *PodContainersGroup) getContainerVolumeMountRelationshipValues(parentExpression, imageID, childClassTypeID string) string {
	c, ok := p.Containers[imageID]
	if !ok {
		return parentExpression
	}
	for _, vm := range c.VolumeMounts {
		parentExpression += fmt.Sprintf("\n('%s', (%s), '%d'),", childClassTypeID, selectRelatedContainerIDFromImageID(imageID), vm.VolumeMountID)
	}
	return parentExpression
}

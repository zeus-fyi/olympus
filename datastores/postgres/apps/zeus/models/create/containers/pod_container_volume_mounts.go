package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

func (p *PodContainersGroup) insertContainerVolumeMounts(parentExpression string, contVolMounts containers.ContainerVolumeMounts) string {
	valsToInsert := "VALUES "

	for i, vm := range contVolMounts {
		valsToInsert += fmt.Sprintf("('%d', '%s', '%s')", vm.VolumeMountID, vm.VolumeMountPath, vm.VolumeMountPath)
		if i < len(contVolMounts)-1 {
			valsToInsert += ","
		}
	}

	containerVmsInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO container_volume_mounts(volume_mount_id, volume_mount_path, volume_name)
					%s
	),`, "cte_container_volume_mounts", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerVmsInsert)
	return returnExpression
}

func (p *PodContainersGroup) insertContainerVolumeMountRelationship(parentExpression string, containerImageID string, vm autogen_structs.ContainerVolumeMounts, cct autogen_structs.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "

	valsToInsert += fmt.Sprintf("('%d', (%s), '%d')", cct.ChartSubcomponentChildClassTypeID, selectRelatedContainerIDFromImageID(containerImageID), vm.VolumeMountID)

	containerVmRelationshipsToInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_volume_mounts(chart_subcomponent_child_class_type_id, container_id, volume_mount_id)
					%s
	),`, "cte_containers_volume_mounts", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerVmRelationshipsToInsert)
	return returnExpression
}

package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

func (p *PodContainersGroup) insertContainerVolumeMountsHeader() string {
	return "INSERT INTO container_probes(probe_id, probe_key_values_jsonb) VALUES "
}

func (p *PodContainersGroup) getInsertContainerVolumeMountsValues(parentExpression string, contVolMounts containers.ContainerVolumeMounts) string {
	for _, vm := range contVolMounts {
		parentExpression += fmt.Sprintf("('%d', '%s', '%s')", vm.VolumeMountID, vm.VolumeMountPath, vm.VolumeMountPath)
	}
	return parentExpression
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

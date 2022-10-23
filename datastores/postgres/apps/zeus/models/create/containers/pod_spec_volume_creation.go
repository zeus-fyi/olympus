package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required

func (p *PodContainersGroup) insertVolumesHeader() string {
	return "INSERT INTO volumes(volume_id, volume_name, volume_key_values_jsonb) "
}

func (p *PodContainersGroup) insertVolumes(parentExpression, childClassTypeID string, vols ...autogen_bases.Volumes) string {
	valsToInsert := "VALUES "

	for i, v := range vols {
		valsToInsert += fmt.Sprintf("('%d', '%s', '%s')", v.VolumeID, v.VolumeName, v.VolumeKeyValuesJSONb)
		if i < len(vols)-1 {
			valsToInsert += ","
		}
	}

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO volumes(volume_id, volume_name, volume_key_values_jsonb)
					%s
	),`, "cte_insert_volumes_with_cv_mount", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

func (p *PodContainersGroup) insertVolumePodSpecRelationship(parentExpression, volID, childClassTypeID string) string {
	valsToInsert := "VALUES "
	valsToInsert += fmt.Sprintf("('%s', %s)", childClassTypeID, volID)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_volumes(chart_subcomponent_child_class_type_id, volume_id)
					%s
	),`, "cte_containers_volumes_spec_relationship", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

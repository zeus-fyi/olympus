package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required
func (p *PodContainersGroup) insertVolumesHeader() string {
	return "INSERT INTO volumes(volume_id, volume_name, volume_key_values_jsonb) "
}

func (p *PodContainersGroup) insertVolumes(cteSubfield sql_query_templates.SubCTE) {
	vols := p.PodSpecTemplate.Spec.PodTemplateSpecVolumes
	for _, v := range vols {
		cteSubfield.AddValues(v.VolumeID, v.VolumeName, v.VolumeKeyValuesJSONb)
	}
	return
}

func (p *PodContainersGroup) insertVolumePodSpecRelationshipsHeader() string {
	return "INSERT INTO containers_volumes(chart_subcomponent_child_class_type_id, volume_id) "
}

func (p *PodContainersGroup) insertVolumePodSpecRelationship(podSpecChildClassTypeID int, cteSubfield sql_query_templates.SubCTE) {
	vols := p.PodSpecTemplate.Spec.PodTemplateSpecVolumes
	for _, v := range vols {
		cteSubfield.AddValues(podSpecChildClassTypeID, v.VolumeID)
	}
	return
}

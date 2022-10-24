package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required
func (p *PodContainersGroup) insertVolumesHeader() string {
	return "INSERT INTO volumes(volume_id, volume_name, volume_key_values_jsonb) "
}

func (p *PodContainersGroup) insertVolumePodSpecRelationshipsHeader() string {
	return "INSERT INTO containers_volumes(chart_subcomponent_child_class_type_id, volume_id) "
}

func (p *PodContainersGroup) insertVolumes(podSpecChildClassTypeID int, podSpecVolumesSubCTE, podSpecVolumesRelationshipSubCTE *sql_query_templates.SubCTE) {
	vols := p.PodSpecTemplate.Spec.PodTemplateSpecVolumes
	ts := chronos.Chronos{}
	for _, v := range vols {
		volID := ts.UnixTimeStampNow()
		podSpecVolumesSubCTE.AddValues(volID, v.VolumeName, v.VolumeKeyValuesJSONb)
		podSpecVolumesRelationshipSubCTE.AddValues(podSpecChildClassTypeID, volID)
	}
}

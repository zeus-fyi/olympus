package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required
func (p *PodTemplateSpec) insertVolumes(podSpecVolumesSubCTE, podSpecVolumesRelationshipSubCTE *sql_query_templates.SubCTE) {
	vols := p.Spec.PodTemplateSpecVolumes
	ts := chronos.Chronos{}
	for _, v := range vols {
		volID := ts.UnixTimeStampNow()
		podSpecVolumesSubCTE.AddValues(volID, v.VolumeName, v.VolumeKeyValuesJSONb)
		podSpecVolumesRelationshipSubCTE.AddValues(p.GetPodSpecChildClassTypeID(), volID)
	}
}

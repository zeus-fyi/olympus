package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required
func (p *PodTemplateSpec) insertVolumes() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// volumes for pod spec
	ts := chronos.Chronos{}

	agVol := autogen_bases.Volumes{}
	podSpecVolumesSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_pod_spec_volumes_%d", ts.UnixTimeStampNow()))
	podSpecVolumesSubCTE.TableName = agVol.GetTableName()
	podSpecVolumesSubCTE.Columns = []string{"volume_id", "volume_name", "volume_key_values_jsonb"}

	agVolR := autogen_bases.ContainersVolumes{}

	podSpecVolumesRelationshipSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_pod_spec_containers_volumes_%d", ts.UnixTimeStampNow()))
	podSpecVolumesRelationshipSubCTE.TableName = agVolR.GetTableName()
	podSpecVolumesRelationshipSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "volume_id"}
	vols := p.Spec.PodTemplateSpecVolumes
	cID := p.GetPodSpecChildClassTypeID()

	for _, v := range vols {
		volID := ts.UnixTimeStampNow()
		podSpecVolumesSubCTE.AddValues(volID, v.VolumeName, v.VolumeKeyValuesJSONb)
		podSpecVolumesRelationshipSubCTE.AddValues(cID, volID)
	}

	return podSpecVolumesSubCTE, podSpecVolumesRelationshipSubCTE
}

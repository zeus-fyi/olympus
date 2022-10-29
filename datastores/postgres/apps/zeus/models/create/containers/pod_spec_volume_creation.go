package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// This will create a volume for the pod spec, if a volume already exists it's not needed, so this is optionally
// required
func (p *PodTemplateSpec) insertVolumes() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// volumes for pod spec
	agVol := autogen_bases.Volumes{}
	podSpecVolumesSubCTE := sql_query_templates.NewSubInsertCTE("cte_pod_spec_volumes")
	podSpecVolumesSubCTE.TableName = agVol.GetTableName()
	podSpecVolumesSubCTE.Columns = agVol.GetTableColumns()

	agVolR := autogen_bases.ContainersVolumes{}

	podSpecVolumesRelationshipSubCTE := sql_query_templates.NewSubInsertCTE("cte_pod_spec_containers_volumes")
	podSpecVolumesRelationshipSubCTE.TableName = agVolR.GetTableName()
	podSpecVolumesRelationshipSubCTE.Columns = agVolR.GetTableColumns()
	vols := p.Spec.PodTemplateSpecVolumes
	ts := chronos.Chronos{}
	cID := p.GetPodSpecChildClassTypeID()

	for _, v := range vols {
		volID := ts.UnixTimeStampNow()
		podSpecVolumesSubCTE.AddValues(volID, v.VolumeName, v.VolumeKeyValuesJSONb)
		podSpecVolumesRelationshipSubCTE.AddValues(cID, volID)
	}

	return podSpecVolumesSubCTE, podSpecVolumesRelationshipSubCTE
}

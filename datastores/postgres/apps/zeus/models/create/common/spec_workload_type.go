package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateSpecWorkloadTypeSubCTE(c *create.Chart, specWorkload structs.SpecWorkload) sql_query_templates.SubCTEs {
	parentClassTypeSubCTE := CreateParentClassTypeSubCTE(&specWorkload.ChartSubcomponentParentClassTypes)
	pcID := specWorkload.ChartSubcomponentParentClassTypeID
	specWorkload.SetParentClassTypeIDs(pcID)
	replicaSubCtes := CreateChildClassSingleValueSubCTEs(&specWorkload.Replicas)
	matchLabelsCtes := CreateChildClassMultiValueSubCTEs(&specWorkload.Selector.MatchLabels)
	chartComponentRelationship := AddParentClassToChartPackage(c, pcID)

	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentClassTypeSubCTE, replicaSubCtes, matchLabelsCtes, []sql_query_templates.SubCTE{chartComponentRelationship})
	return combinedSubCtes
}

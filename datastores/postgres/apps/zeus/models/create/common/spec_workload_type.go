package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateSpecWorkloadTypeSubCTE(specWorkload common.SpecWorkload) sql_query_templates.SubCTEs {
	if specWorkload.ChartSubcomponentParentClassTypeID == 0 {
		var ts chronos.Chronos
		pcTypeClassTypeID := ts.UnixTimeStampNow()
		specWorkload.ChartSubcomponentParentClassTypeID = pcTypeClassTypeID
	}

	parentClassTypeSubCTE := CreateParentClassTypeSubCTE(specWorkload.ChartSubcomponentParentClassTypes)
	replicaSubCtes := CreateChildClassSingleValueSubCTEs(specWorkload.Replicas)
	matchLabelsCtes := CreateChildClassMultiValueSubCTEs(specWorkload.Selector.MatchLabels)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentClassTypeSubCTE, replicaSubCtes, matchLabelsCtes)
	return combinedSubCtes
}

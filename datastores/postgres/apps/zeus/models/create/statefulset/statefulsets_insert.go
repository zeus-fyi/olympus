package statefulset

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulset"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type StatefulSet struct {
	statefulset.StatefulSet
}

const ModelName = "StatefulSet"

func (s *StatefulSet) InsertStatefulSet(ctx context.Context, q sql_query_templates.QueryParams, c create.Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	q.CTEQuery = s.InsertStatefulSetCte()
	q.RawQuery = q.CTEQuery.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

func (s *StatefulSet) InsertStatefulSetCte() sql_query_templates.CTE {
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(s.Metadata)
	// spec
	specCtes := common.CreateSpecWorkloadTypeSubCTE(s.Spec.SpecWorkload)

	// pod template spec
	podSpecTemplateCte := s.Spec.Template.InsertPodTemplateSpecContainersCTE()
	podSpecTemplateSubCtes := podSpecTemplateCte.SubCTEs

	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes, podSpecTemplateSubCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertStatefulSetCte",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

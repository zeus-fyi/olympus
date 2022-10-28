package networking

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Service struct {
	networking.Service
}

const ModelName = "Service"

func (s *Service) InsertService(ctx context.Context, q sql_query_templates.QueryParams, c *charts.Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	q.CTEQuery = s.InsertServiceCte(c)
	q.RawQuery = q.CTEQuery.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertService: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

func (s *Service) InsertServiceCte(chart *charts.Chart) sql_query_templates.CTE {

	if chart != nil {
		s.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)
	// spec
	specCtes := s.CreateServiceSpecSubCTE(chart)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertServiceCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

package packages

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Packages struct {
	charts.Chart
	CTE sql_query_templates.CTE
	chart_workload.ChartWorkload
}

func NewPackageInsert() Packages {
	cw := chart_workload.NewChartWorkload()

	pkg := Packages{
		Chart:         charts.NewChart(),
		CTE:           sql_query_templates.CTE{},
		ChartWorkload: cw,
	}
	return pkg
}

const Sn = "Packages"

func (p *Packages) InsertPackages(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	q.CTEQuery = p.InsertPackagesCTE()
	q.RawQuery = q.CTEQuery.GenerateChainedCTE()

	r, err := apps.Pg.Exec(ctx, q.RawQuery, q.CTEQuery.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Packages: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (p *Packages) InsertPackagesCTE() sql_query_templates.CTE {
	if p.Deployment != nil {
		depCte := p.GetDeploymentCTE(&p.Chart)
		p.CTE.AppendSubCtes(depCte.SubCTEs)
	}
	if p.StatefulSet != nil {
		stsCte := p.GetStatefulSetCTE(&p.Chart)
		p.CTE.AppendSubCtes(stsCte.SubCTEs)
	}
	if p.Service != nil {
		svcCte := p.GetServiceCTE(&p.Chart)
		p.CTE.AppendSubCtes(svcCte.SubCTEs)
	}
	if p.Ingress != nil {
		ingressCte := p.GetIngressCTE(&p.Chart)
		p.CTE.AppendSubCtes(ingressCte.SubCTEs)
	}
	if p.ConfigMap != nil {
		cmCte := p.GetConfigMapCTE(&p.Chart)
		p.CTE.AppendSubCtes(cmCte.SubCTEs)
	}
	return p.CTE
}

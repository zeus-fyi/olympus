package packages

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Packages struct {
	charts.Chart
	*deployments.Deployment
	*services.Service
	*ingresses.Ingress
	*configuration.ConfigMap
}

const Sn = "Packages"

func (p *Packages) InsertPackages(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	q.CTEQuery = p.InsertPackagesCTE()
	q.RawQuery = q.CTEQuery.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Packages: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (p *Packages) InsertPackagesCTE() sql_query_templates.CTE {
	var cte sql_query_templates.CTE
	if p.Deployment != nil {
		depCte := p.GetDeploymentCTE(&p.Chart)
		cte.AppendSubCtes(depCte.SubCTEs)
	}
	if p.Service != nil {
		svcCte := p.GetServiceCTE(&p.Chart)
		cte.AppendSubCtes(svcCte.SubCTEs)
	}
	if p.Ingress != nil {
		ingressCte := p.GetIngressCTE(&p.Chart)
		cte.AppendSubCtes(ingressCte.SubCTEs)
	}
	if p.ConfigMap != nil {
		cmCte := p.GetConfigMapCTE(&p.Chart)
		cte.AppendSubCtes(cmCte.SubCTEs)
	}
	return cte
}

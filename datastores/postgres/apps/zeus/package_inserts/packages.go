package package_inserts

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Packages struct {
	charts.Chart
	deployments.Deployment
	networking.Service
}

const Sn = "Packages"

func (p *Packages) InsertPackages(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	q.CTEQuery = p.InsertPackagesCTE()

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

	return cte
}

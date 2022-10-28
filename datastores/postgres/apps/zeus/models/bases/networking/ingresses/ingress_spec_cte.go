package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (i *Ingress) GetIngressSpecCTE(chart *charts.Chart) sql_query_templates.SubCTEs {

	return sql_query_templates.SubCTEs{}
}

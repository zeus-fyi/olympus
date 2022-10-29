package secrets

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *Secret) GetSecretCTE(c *charts.Chart) sql_query_templates.CTE {
	return sql_query_templates.CTE{}
}

package networking

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Service struct {
	networking.Service
}

func (s *Service) InsertService(ctx context.Context, q sql_query_templates.QueryParams, c *create.Chart) error {
	return nil
}

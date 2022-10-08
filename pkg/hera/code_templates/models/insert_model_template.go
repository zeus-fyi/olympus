package models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *StructNameExamples) InsertStructNameExamplesFieldCase(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(sn))
	r, err := postgres_apps.Pg.Exec(ctx, q.SelectQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(sn))
}

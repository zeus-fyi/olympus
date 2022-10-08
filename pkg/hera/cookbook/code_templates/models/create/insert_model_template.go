package create

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/code_templates/models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type StructNameExamples struct {
	models.StructNameExamples
}

func (s *StructNameExamples) StructNameExamplesFieldCase(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(models.Sn))
	r, err := apps.Pg.Exec(ctx, q.SelectQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(models.Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(models.Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(models.Sn))
}

package read

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/hera/code_templates/models"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type StructNameExamples struct {
	models.StructNameExamples
}

func (s *StructNameExamples) StructNameExamplesFieldCase(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(models.Sn))
	rows, err := apps.Pg.Query(ctx, q.SelectQuery())
	if err != nil {
		log.Err(err).Msg(q.LogHeader(models.Sn))
		return err
	}
	defer rows.Close()
	var selectedStructNameExamples models.StructNameExamples
	for rows.Next() {
		var se models.StructNameExample
		rowErr := rows.Scan(se.GetRowValues(q.QueryName))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(models.Sn))
			return rowErr
		}
		selectedStructNameExamples = append(selectedStructNameExamples, se)
	}
	return nil
}

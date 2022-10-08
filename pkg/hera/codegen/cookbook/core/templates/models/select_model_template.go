package models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *StructNameExamples) SelectStructNameExamplesFieldCase(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(sn))
	rows, err := postgres_apps.Pg.Query(ctx, q.SelectQuery())
	if err != nil {
		log.Err(err).Msg(q.LogHeader(sn))
		return err
	}
	defer rows.Close()
	var selectedStructNameExamples StructNameExamples
	for rows.Next() {
		var se StructNameExample
		rowErr := rows.Scan(se.GetRowValues(q.QueryName))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(sn))
			return rowErr
		}
		selectedStructNameExamples = append(selectedStructNameExamples, se)
	}
	return nil
}

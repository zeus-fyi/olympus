package deployments

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/deployments"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Deployment struct {
	deployments.Deployment
}

const ModelName = "Deployment"

func (s *Deployment) SelectDeploymentResource(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.SelectQuery())
	if err != nil {
		log.Err(err).Msg(q.LogHeader(ModelName))
		return err
	}
	defer rows.Close()
	//var podTemplateSpec containers.PodTemplateSpec
	for rows.Next() {
		//var se models.StructNameExample

		rowErr := rows.Scan()
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return rowErr
		}
		//selectedStructNameExamples = append(selectedStructNameExamples, se)
	}
	return nil
}

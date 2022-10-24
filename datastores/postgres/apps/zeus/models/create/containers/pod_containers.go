package containers

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) NewPodContainersMapForDB() map[string]containers.Container {
	m := make(map[string]containers.Container)
	for _, c := range p.GetContainers() {
		m[c.Metadata.ContainerImageID] = c
	}
	return m
}

const ModelName = "PodContainersGroup"

func (p *PodTemplateSpec) InsertPodTemplateSpec(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	q.RawQuery = p.InsertPodContainerGroupSQL()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

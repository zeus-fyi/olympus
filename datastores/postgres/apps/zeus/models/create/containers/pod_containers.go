package containers

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type PodContainersGroup struct {
	PodSpecTemplate containers.PodTemplateSpec
	Containers      map[string]containers.Container
}

func NewPodContainersGroupForDB(ps containers.PodTemplateSpec) PodContainersGroup {
	m := make(map[string]containers.Container)
	for _, c := range ps.GetContainers() {
		m[c.Metadata.ContainerImageID] = c
	}
	return PodContainersGroup{PodSpecTemplate: ps, Containers: m}
}

const ModelName = "PodContainersGroup"

func (p *PodContainersGroup) InsertPodContainerGroup(ctx context.Context, q sql_query_templates.QueryParams, podSpecClassTypeID int) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(ModelName))
	query := p.insertPodContainerGroupSQL(podSpecClassTypeID)
	// TODO create real query
	r, err := apps.Pg.Exec(ctx, query.Name)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("StructNameExamples: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

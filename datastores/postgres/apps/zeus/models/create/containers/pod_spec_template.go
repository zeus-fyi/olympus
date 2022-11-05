package containers

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type PodSpecContainerMetadata struct {
	autogen_bases.ChartSubcomponentSpecPodTemplateContainers
}

var PsName = "ChartSubcomponentSpecPodTemplateContainers"

func (p *PodSpecContainerMetadata) insertChartSubcomponentSpecPodTemplateContainers() string {
	columns := p.GetTableColumns()
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO %s(%s)
 				 VALUES ('%d', '%d', %d)`,
		p.GetTableName(), strings.Join(columns, ","), p.ContainerSortOrder, p.ChartSubcomponentChildClassTypeID, p.ContainerID)
	return sqlInsertStatement
}

func (p *PodSpecContainerMetadata) InsertChartSubcomponentSpecPodTemplateContainers(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(PsName))
	query := p.insertChartSubcomponentSpecPodTemplateContainers()
	_, err := apps.Pg.Exec(ctx, query)
	return misc.ReturnIfErr(err, q.LogHeader(PsName))
}

package containers

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/common"
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
 				 VALUES ('%t', '%d', '%d', %d)`,
		p.GetTableName(), strings.Join(columns, ","), p.IsInitContainer, p.ContainerSortOrder, p.ChartSubcomponentChildClassTypeID, p.ContainerID)
	return sqlInsertStatement
}

func (p *PodSpecContainerMetadata) InsertChartSubcomponentSpecPodTemplateContainers(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(PsName))
	query := p.insertChartSubcomponentSpecPodTemplateContainers()
	_, err := apps.Pg.Exec(ctx, query)
	return misc.ReturnIfErr(err, q.LogHeader(PsName))
}

func SetPodSpecTemplateChildTypeInsert(parentExpression, parentClassTypeCteName string, psChildClassType autogen_bases.ChartSubcomponentChildClassTypes) string {
	cvTypeName := psChildClassType.ChartSubcomponentChildClassTypeName
	typeInsert, cvTypeNameCte := common.SetCvTypeInsert(cvTypeName, parentClassTypeCteName)

	// TODO needs to link containers into here, so probably even create the containers completely independently
	cvChildCteName := fmt.Sprintf("cte_ps_cvs_%s", cvTypeName)
	valueInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponent_spec_pod_template_containers(chart_subcomponent_child_class_type_id, container_id, is_init_container)
					VALUES ((SELECT chart_subcomponent_child_class_type_id FROM %s), '%s', '%s')
	),`, cvChildCteName, cvTypeNameCte, "", "")

	returnExpression := fmt.Sprintf("%s %s %s", parentExpression, typeInsert, valueInsert)
	return returnExpression
}

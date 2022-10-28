package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

var Sn = "ChartSubcomponentParentClassTypes"

type ParentClass struct {
	autogen_bases.ChartSubcomponentParentClassTypes
}

func (p *ParentClass) InsertParentClassTypeID(id int) {
	p.ChartSubcomponentParentClassTypeID = id
}

const SelectDeploymentResourceID = "(SELECT chart_component_resource_id FROM chart_component_resources WHERE chart_component_kind_name = 'Deployment' AND chart_component_api_version = 'apps/v1')"

func CreateParentClassTypeSubCTE(c *create.Chart, pcType *autogen_bases.ChartSubcomponentParentClassTypes) sql_query_templates.SubCTEs {
	if pcType.ChartSubcomponentParentClassTypeID == 0 {
		var ts chronos.Chronos
		pcTypeClassTypeID := ts.UnixTimeStampNow()
		pcType.ChartSubcomponentParentClassTypeID = pcTypeClassTypeID
	}
	pcType.ChartPackageID = c.GetChartPackageID()
	parentClassTypeSubCTE := createParentClassTypeSubCTE(pcType)
	return []sql_query_templates.SubCTE{parentClassTypeSubCTE}
}

func createParentClassTypeSubCTE(pcType *autogen_bases.ChartSubcomponentParentClassTypes) sql_query_templates.SubCTE {
	queryName := "cte_" + pcType.ChartSubcomponentParentClassTypeName
	parentClassTypeSubCTE := sql_query_templates.NewSubInsertCTE(queryName)
	parentClassTypeSubCTE.TableName = pcType.GetTableName()
	parentClassTypeSubCTE.Fields = pcType.GetTableColumns()
	parentClassTypeSubCTE.Values = []apps.RowValues{pcType.GetRowValues(queryName)}
	return parentClassTypeSubCTE
}

func (p *ParentClass) insertChartSubcomponentParentClassType() string {
	columns := p.GetTableColumns()
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO %s(%s)
 				 VALUES ('%d', '%d', '%d', '%s')`,
		p.GetTableName(), strings.Join(columns, ","), p.ChartPackageID, p.ChartComponentResourceID, p.ChartSubcomponentParentClassTypeID, p.ChartSubcomponentParentClassTypeName)
	return sqlInsertStatement
}

func (p *ParentClass) InsertChartSubcomponentParentClassTypes(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	query := p.insertChartSubcomponentParentClassType()
	_, err := apps.Pg.Exec(ctx, query)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

package common

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

var Sn = "ChartSubcomponentParentClassTypes"

type ParentClass struct {
	autogen_bases.ChartSubcomponentParentClassTypes
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

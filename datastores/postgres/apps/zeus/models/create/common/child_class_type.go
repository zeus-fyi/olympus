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

var StructName = "ChartSubcomponentChildClassTypes"

type ChildClass struct {
	autogen_bases.ChartSubcomponentChildClassTypes
}

func (c *ChildClass) insertChartSubcomponentChildClassTypes() string {
	columns := c.GetTableColumns()
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO %s(%s)
 				 VALUES ('%d', '%d', '%s')`,
		c.GetTableName(), strings.Join(columns, ","), c.ChartSubcomponentParentClassTypeID, c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChildClassTypeName)
	return sqlInsertStatement
}

func (c *ChildClass) InsertChartSubcomponentChildClassTypes(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(StructName))
	query := c.insertChartSubcomponentChildClassTypes()
	_, err := apps.Pg.Exec(ctx, query)
	return misc.ReturnIfErr(err, q.LogHeader(StructName))
}

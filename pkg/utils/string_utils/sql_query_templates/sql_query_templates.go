package sql_query_templates

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type Query struct{}

type QueryParams struct {
	QueryName   string
	RawQuery    string
	CTEQuery    CTE
	Columns     []string
	TableName   string
	WhereClause string
	Values      []apps.RowValues
	Limit       int
}

// NewQueryParam queryName, tableName, whereClause string, limit int, fields []string
func NewQueryParam(queryName, tableName, whereClause string, limit int, columns []string) QueryParams {
	return QueryParams{
		QueryName:   queryName,
		Columns:     columns,
		TableName:   tableName,
		WhereClause: whereClause,
		Limit:       limit,
	}
}

func (q *QueryParams) LogHeader(structName string) string {
	query := fmt.Sprintf(`%s: QueryName: %s, TableName: %s, WhereClause %s, Limit %d`, structName, q.Columns, q.TableName, q.WhereClause, q.Limit)
	return query
}

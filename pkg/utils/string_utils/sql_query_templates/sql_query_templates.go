package sql_query_templates

import "fmt"

type Query struct{}

type QueryParams struct {
	QueryName   string
	Fields      []string
	TableName   string
	WhereClause string
	Values      []any
	Limit       int
}

// NewQueryParam queryName, tableName, whereClause string, limit int, fields []string
func NewQueryParam(queryName, tableName, whereClause string, limit int, fields []string) QueryParams {
	return QueryParams{
		QueryName:   queryName,
		Fields:      fields,
		TableName:   tableName,
		WhereClause: whereClause,
		Limit:       limit,
	}
}

func (q *QueryParams) LogHeader(structName string) string {
	query := fmt.Sprintf(`%s: QueryName: %s, TableName: %s, WhereClause %s, Limit %d`, structName, q.Fields, q.TableName, q.WhereClause, q.Limit)
	return query
}

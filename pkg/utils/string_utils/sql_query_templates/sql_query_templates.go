package sql_query_templates

import "fmt"

type QueryParams struct {
	QueryName   string
	Fields      []string
	TableName   string
	WhereClause string
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

func (q *QueryParams) SelectQuery() string {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s LIMIT %d`, q.Fields, q.TableName, q.WhereClause, q.Limit)
	return query
}

func (q *QueryParams) LogHeader(structName string) string {
	query := fmt.Sprintf(`%s: QueryName: %s, TableName: %s, WhereClause %s, Limit %d`, structName, q.Fields, q.TableName, q.WhereClause, q.Limit)
	return query
}

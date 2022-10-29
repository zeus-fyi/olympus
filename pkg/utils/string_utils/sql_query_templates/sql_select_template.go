package sql_query_templates

import "fmt"

func (q *QueryParams) SelectQuery() string {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s LIMIT %d`, q.Columns, q.TableName, q.WhereClause, q.Limit)
	return query
}

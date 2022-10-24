package sql_query_templates

import (
	"fmt"
	"strings"
)

func (q *QueryParams) InsertCTEHeader() string {
	query := fmt.Sprintf("%s AS (\n INSERT INTO %s(%s) VALUES \n", q.QueryName, q.TableName, strings.Join(q.Fields, ","))
	return query
}

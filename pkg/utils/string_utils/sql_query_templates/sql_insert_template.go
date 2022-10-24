package sql_query_templates

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (q *QueryParams) InsertQueryHeader() string {
	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES `, q.TableName, q.Fields)
	return query
}

func (q *QueryParams) AddValues(values ...any) {
	tmp := make(apps.RowValues, len(values))
	for i, v := range values {
		tmp[i] = v
	}
	q.Values = append(q.Values, tmp)
	return
}

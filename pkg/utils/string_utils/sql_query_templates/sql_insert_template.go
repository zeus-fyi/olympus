package sql_query_templates

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (q *QueryParams) InsertSingleElementQuery() string {
	vals := string_utils.StringDelimitedSliceBuilderSQLMultiRowValues(",", q.Values)
	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES %s`, q.TableName, strings.Join(q.Columns, ","), vals)
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

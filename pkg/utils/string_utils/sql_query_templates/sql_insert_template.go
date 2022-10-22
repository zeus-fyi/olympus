package sql_query_templates

import "fmt"

func (q *QueryParams) InsertQueryHeader() string {
	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES `, q.TableName, q.Fields)
	return query
}

func (q *QueryParams) AddValues(values ...any) {
	if len(q.Values) <= 0 {
		q.Values = make([]any, len(values))
	}
	for _, v := range values {
		q.Values = append(q.Values, v)
	}
}

package models

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateTestQueryNameParams() sql_query_templates.QueryParams {
	queryName := "queryName"
	tableName := "tableName"
	param := "param"
	whereClause := fmt.Sprintf("fieldName = %s", param)
	limit := 10000
	fields := []string{"fieldName1", "fieldName2"}

	return sql_query_templates.NewQueryParam(queryName, tableName, whereClause, limit, fields)
}

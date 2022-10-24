package sql_query_templates

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *SubCTE) GenerateChainedInsertCTE() string {
	formattedValues := MultiLevelValuesCTEStringBuilderSQL(s.InsertCTEHeader(), s.Values)
	return formattedValues
}

func MultiLevelValuesCTEStringBuilderSQL(insertHeader string, values []apps.RowValues) string {
	sb := strings.Builder{}
	sb.WriteString(insertHeader)

	for i, line := range values {
		sb.WriteString("(")

		for c, val := range line {
			switch val.(type) {
			case string:
				sb.WriteString(val.(string))
			case int:
				intStr := fmt.Sprintf("%d", val.(int))
				sb.WriteString(intStr)
			case int64:
				int64Str := fmt.Sprintf("%d", val.(int64))
				sb.WriteString(int64Str)
			case uint64:
				uint64Str := fmt.Sprintf("%d", val.(uint64))
				sb.WriteString(uint64Str)
			case bool:
				boolStr := fmt.Sprintf("%t", val.(bool))
				sb.WriteString(boolStr)
			default:
			}
			if len(line)-1 != c {
				sb.WriteString(", ")
			}
		}

		if len(values)-1 == i {
			sb.WriteString(")\n")
			return sb.String()
		}
		sb.WriteString("),\n")
	}
	return sb.String()
}

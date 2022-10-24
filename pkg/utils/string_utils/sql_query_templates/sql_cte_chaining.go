package sql_query_templates

import (
	"fmt"
	"strings"
)

func (c *CTE) MultiLevelValuesCTEStringBuilderSQL() string {
	sb := strings.Builder{}
	sb.WriteString("WITH ")

	ctesWithValues := 0
	for _, se := range c.SubCTEs {
		if len(se.Values) > 0 {
			ctesWithValues += 1
		}
	}

	for subCteCount, subCteExpr := range c.SubCTEs {

		if len(subCteExpr.Values) > 0 {
			sb.WriteString(subCteExpr.InsertCTEHeader())

			for i, line := range subCteExpr.Values {
				sb.WriteString("(")
				for col, val := range line {
					switch val.(type) {
					case string:
						stringField := val.(string)
						if strings.HasPrefix(stringField, "(SELECT") {
							sb.WriteString(val.(string))
						} else {
							sb.WriteString("'")
							sb.WriteString(val.(string))
							sb.WriteString("'")
						}
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
					if len(line)-1 != col {
						sb.WriteString(", ")
					}
				}
				if len(subCteExpr.Values)-1 == i {
					sb.WriteString(")")
				} else {
					sb.WriteString("),")
				}
			}
		}

		if ctesWithValues-1 == subCteCount {
			sb.WriteString(")\n SELECT 1")
			return sb.String()
		} else {
			sb.WriteString("\n),\n")
		}
	}
	return sb.String()
}

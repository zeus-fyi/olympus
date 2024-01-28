package sql_query_templates

import (
	"fmt"
	"strings"
)

func (c *CTE) SanitizedMultiLevelValuesCTEStringBuilderSQL() string {
	sb := strings.Builder{}
	sb.WriteString("WITH ")

	var filteredSubCTEs SubCTEs
	for _, se := range c.SubCTEs {
		if len(se.Values) > 0 || len(se.RawQuery) > 0 {
			filteredSubCTEs = append(filteredSubCTEs, se)
		}
	}
	for cteCount, subCteExpr := range filteredSubCTEs {

		if len(subCteExpr.RawQuery) <= 0 {
			sb.WriteString(subCteExpr.InsertCTEHeader())
			for lc, line := range subCteExpr.Values {
				sb.WriteString("(")
				for col, val := range line {
					if v, ok := subCteExpr.ValuesOverride[col]; ok {
						sb.WriteString(v)
					} else {
						c.Params = append(c.Params, val)
						sb.WriteString(fmt.Sprintf("$%d", len(c.Params)))
					}
					if len(line)-1 != col {
						sb.WriteString(", ")
					}
				}
				if len(subCteExpr.Values)-1 == lc {
					sb.WriteString(")")
				} else {
					sb.WriteString("), \n")
				}
			}
			if len(c.OnConflicts) > 0 {
				sb.WriteString(fmt.Sprintf("\n ON CONFLICT (%s) DO UPDATE SET ", strings.Join(c.OnConflicts, ", ")))
				for i, col := range c.OnConflictsUpdateColumns {
					sb.WriteString(fmt.Sprintf("%s = EXCLUDED.%s", col, col))
					if i < len(c.OnConflictsUpdateColumns)-1 {
						sb.WriteString(", ")
					}
				}
			}
		} else {
			sb.WriteString(subCteExpr.InjectRawQueryCTE())
		}

		if len(filteredSubCTEs)-1 == cteCount {
			if len(c.ReturnSQLStatement) > 0 {
				sb.WriteString(")\n" + c.ReturnSQLStatement)
			} else {
				if c.OnConflictDoNothing {
					sb.WriteString("\n ON CONFLICT DO NOTHING)\n SELECT 1")
				} else {
					sb.WriteString(")\n SELECT 1")
				}
			}
			return sb.String()
		} else {
			sb.WriteString("\n),\n")
		}
	}
	return sb.String()
}

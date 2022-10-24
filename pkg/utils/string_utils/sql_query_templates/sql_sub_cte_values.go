package sql_query_templates

import "github.com/zeus-fyi/olympus/pkg/utils/string_utils"

func (s *SubCTE) GenerateChainedInsertCTE() string {
	tmp := ""
	for i, expr := range s.Values {
		if len(s.Values)-1 == i {
			tmp = string_utils.StringDelimitedSliceBuilderSQL(",", expr)
			return tmp
		} else {
			tmp = string_utils.StringDelimitedSliceBuilderSQL(",", expr)
		}
	}
	return tmp
}

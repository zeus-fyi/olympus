package common_fields

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"

func QueryParams() fields.Field {
	return fields.Field{
		Pkg:   "sql_query_templates",
		Name:  "q",
		Type:  "QueryParams",
		Value: "",
	}
}

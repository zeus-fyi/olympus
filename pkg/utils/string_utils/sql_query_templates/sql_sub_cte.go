package sql_query_templates

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type SubCTEs []SubCTE

type SubCTE struct {
	QueryParams
}

func NewSubInsertCTE(subCteName string) SubCTE {
	return SubCTE{
		QueryParams{
			QueryName:   subCteName,
			Fields:      nil,
			TableName:   "",
			WhereClause: "",
			Values:      []apps.RowValues{},
			Limit:       0,
		},
	}
}

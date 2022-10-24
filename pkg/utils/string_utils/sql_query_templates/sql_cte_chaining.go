package sql_query_templates

type CTE struct {
	Name string
	SubCTEs
}

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
			Values:      [][]any{},
			Limit:       0,
		},
	}
}

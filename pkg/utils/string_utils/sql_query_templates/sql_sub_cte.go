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
			Columns:     nil,
			TableName:   "",
			WhereClause: "",
			Values:      []apps.RowValues{},
			Limit:       0,
		},
	}
}

func AppendSubCteSlices(singleSubCte SubCTEs, subCteGroup ...SubCTEs) SubCTEs {
	combinedSubCtes := singleSubCte
	for _, addSubCte := range subCteGroup {
		combinedSubCtes = append(combinedSubCtes, addSubCte...)
	}
	return combinedSubCtes
}

package sql_query_templates

type CTE struct {
	Name string
	SubCTEs
	Params             []interface{}
	ReturnSQLStatement string
}

func (c *CTE) GenerateChainedCTE() string {
	formattedValues := c.SanitizedMultiLevelValuesCTEStringBuilderSQL()
	return formattedValues
}

func (c *CTE) AppendSubCtes(se SubCTEs) {
	c.SubCTEs = append(c.SubCTEs, se...)
}

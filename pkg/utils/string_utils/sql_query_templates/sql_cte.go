package sql_query_templates

type CTE struct {
	Name string
	SubCTEs
	ReturnSQLStatement string
}

func (c *CTE) GenerateChainedCTE() string {
	formattedValues := c.MultiLevelValuesCTEStringBuilderSQL()
	return formattedValues
}

func (c *CTE) AppendSubCtes(se SubCTEs) {
	c.SubCTEs = append(c.SubCTEs, se...)
}

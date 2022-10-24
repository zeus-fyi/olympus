package sql_query_templates

type CTE struct {
	Name string
	SubCTEs
}

func (c *CTE) GenerateChainedCTE() string {
	formattedValues := c.MultiLevelValuesCTEStringBuilderSQL()
	return formattedValues
}

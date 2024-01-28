package sql_query_templates

type CTE struct {
	OnConflicts              []string
	OnConflictsUpdateColumns []string
	Name                     string
	SubCTEs
	Params              []interface{} // FYI this stupid thing will show null values using goland IDE
	ReturnSQLStatement  string
	OnConflictDoNothing bool
}

func (c *CTE) GenerateChainedCTE() string {
	formattedValues := c.SanitizedMultiLevelValuesCTEStringBuilderSQL()
	return formattedValues
}

func (c *CTE) AppendSubCtes(se SubCTEs) {
	c.SubCTEs = append(c.SubCTEs, se...)
}

package autogen_bases

type TopologyInfrastructureComponents struct {
	ChartPackageID int `db:"chart_package_id"`
	TopologyID     int `db:"topology_id"`
}
type TopologyInfrastructureComponentsSlice []TopologyInfrastructureComponents

func (t *TopologyInfrastructureComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.ChartPackageID, t.TopologyID}
	}
	return pgValues
}
func (t *TopologyInfrastructureComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_package_id", "topology_id"}
	return columnValues
}
func (t *TopologyInfrastructureComponents) GetTableName() (tableName string) {
	tableName = "topology_infrastructure_components"
	return tableName
}

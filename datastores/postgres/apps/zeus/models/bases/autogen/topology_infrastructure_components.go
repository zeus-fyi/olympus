package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyInfrastructureComponents struct {
	TopologyID     int `db:"topology_id"`
	ChartPackageID int `db:"chart_package_id"`
}
type TopologyInfrastructureComponentsSlice []TopologyInfrastructureComponents

func (t *TopologyInfrastructureComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.ChartPackageID}
	}
	return pgValues
}
func (t *TopologyInfrastructureComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "chart_package_id"}
	return columnValues
}
func (t *TopologyInfrastructureComponents) GetTableName() (tableName string) {
	tableName = "topology_infrastructure_components"
	return tableName
}

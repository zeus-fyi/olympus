package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyInfrastructureComponents struct {
	TopologyID                        int `db:"topology_id" json:"topologyID"`
	ChartPackageID                    int `db:"chart_package_id" json:"chartPackageID"`
	TopologyInfrastructureComponentID int `db:"topology_infrastructure_component_id" json:"topologyInfrastructureComponentID"`
}
type TopologyInfrastructureComponentsSlice []TopologyInfrastructureComponents

func (t *TopologyInfrastructureComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.ChartPackageID, t.TopologyInfrastructureComponentID}
	}
	return pgValues
}
func (t *TopologyInfrastructureComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "chart_package_id", "topology_infrastructure_component_id"}
	return columnValues
}
func (t *TopologyInfrastructureComponents) GetTableName() (tableName string) {
	tableName = "topology_infrastructure_components"
	return tableName
}

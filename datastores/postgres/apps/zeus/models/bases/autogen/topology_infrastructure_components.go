package autogen_bases

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type TopologyInfrastructureComponents struct {
	TopologyInfrastructureComponentID int `db:"topology_infrastructure_component_id" json:"topologyInfrastructureComponentID"`
	TopologyID                        int `db:"topology_id" json:"topologyID"`
	ChartPackageID                    int `db:"chart_package_id" json:"chartPackageID"`
	TopologySkeletonBaseVersionID     int `db:"topology_skeleton_base_version_id" json:"topologySkeletonBaseVersionID"`
}
type TopologyInfrastructureComponentsSlice []TopologyInfrastructureComponents

func (t *TopologyInfrastructureComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyInfrastructureComponentID, t.TopologyID, t.ChartPackageID, t.TopologySkeletonBaseVersionID}
	}
	return pgValues
}
func (t *TopologyInfrastructureComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_infrastructure_component_id", "topology_id", "chart_package_id", "topology_skeleton_base_version_id"}
	return columnValues
}
func (t *TopologyInfrastructureComponents) GetTableName() (tableName string) {
	tableName = "topology_infrastructure_components"
	return tableName
}

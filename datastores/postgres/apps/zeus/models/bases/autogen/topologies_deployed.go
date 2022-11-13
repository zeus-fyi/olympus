package autogen_bases

type TopologiesDeployed struct {
	DeploymentID   int       `db:"deployment_id" json:"deploymentID"`
	TopologyID     int       `db:"topology_id" json:"topologyID"`
	TopologyStatus string    `db:"topology_status" json:"topologyStatus"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
type TopologiesDeployedSlice []TopologiesDeployed

func (t *TopologiesDeployed) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.DeploymentID, t.TopologyID, t.TopologyStatus, t.UpdatedAt}
	}
	return pgValues
}
func (t *TopologiesDeployed) GetTableColumns() (columnValues []string) {
	columnValues = []string{"deployment_id", "topology_id", "topology_status", "updated_at"}
	return columnValues
}
func (t *TopologiesDeployed) GetTableName() (tableName string) {
	tableName = "topologies_deployed"
	return tableName
}

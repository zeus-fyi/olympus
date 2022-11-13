package autogen_bases

type TopologiesKns struct {
	Region        string `db:"region" json:"region"`
	Context       string `db:"context" json:"context"`
	Namespace     string `db:"namespace" json:"namespace"`
	Env           string `db:"env" json:"env"`
	TopologyID    int    `db:"topology_id" json:"topologyID"`
	CloudProvIDer string `db:"cloud_provider" json:"cloudProvider"`
}
type TopologiesKnsSlice []TopologiesKns

func (t *TopologiesKns) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.Region, t.Context, t.Namespace, t.Env, t.TopologyID, t.CloudProvIDer}
	}
	return pgValues
}
func (t *TopologiesKns) GetTableColumns() (columnValues []string) {
	columnValues = []string{"region", "context", "namespace", "env", "topology_id", "cloud_provider"}
	return columnValues
}
func (t *TopologiesKns) GetTableName() (tableName string) {
	tableName = "topologies_kns"
	return tableName
}

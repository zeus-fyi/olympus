package autogen_bases

type TopologiesKns struct {
	Namespace  string `db:"namespace"`
	Env        string `db:"env"`
	TopologyID int    `db:"topology_id"`
	Context    string `db:"context"`
}
type TopologiesKnsSlice []TopologiesKns

func (t *TopologiesKns) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.Namespace, t.Env, t.TopologyID, t.Context}
	}
	return pgValues
}
func (t *TopologiesKns) GetTableColumns() (columnValues []string) {
	columnValues = []string{"namespace", "env", "topology_id", "context"}
	return columnValues
}
func (t *TopologiesKns) GetTableName() (tableName string) {
	tableName = "topologies_kns"
	return tableName
}

package autogen_bases

type TopologiesKns struct {
	TopologyID int    `db:"topology_id"`
	Context    string `db:"context"`
	Namespace  string `db:"namespace"`
	Env        string `db:"env"`
}
type TopologiesKnsSlice []TopologiesKns

func (t *TopologiesKns) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.Context, t.Namespace, t.Env}
	}
	return pgValues
}
func (t *TopologiesKns) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "context", "namespace", "env"}
	return columnValues
}
func (t *TopologiesKns) GetTableName() (tableName string) {
	tableName = "topologies_kns"
	return tableName
}

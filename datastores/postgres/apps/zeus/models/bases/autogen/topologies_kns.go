package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologiesKns struct {
	TopologyID    int    `db:"topology_id" json:"topology_id"`
	CloudProvider string `db:"cloud_provider" json:"cloud_provider"`
	Region        string `db:"region" json:"region"`
	Context       string `db:"context" json:"context"`
	Namespace     string `db:"namespace" json:"namespace"`
	Env           string `db:"env" json:"env"`
}
type TopologiesKnsSlice []TopologiesKns

func (t *TopologiesKns) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.CloudProvider, t.Region, t.Context, t.Namespace, t.Env}
	}
	return pgValues
}
func (t *TopologiesKns) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "cloud_provider", "region", "context", "namespace", "env"}
	return columnValues
}
func (t *TopologiesKns) GetTableName() (tableName string) {
	tableName = "topologies_kns"
	return tableName
}

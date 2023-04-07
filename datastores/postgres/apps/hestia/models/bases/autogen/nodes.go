package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Nodes struct {
	NodeID        int     `db:"node_id" json:"nodeID"`
	Description   string  `db:"description" json:"description"`
	Slug          string  `db:"slug" json:"slug"`
	Disk          int     `db:"disk" json:"disk"`
	PriceHourly   float64 `db:"price_hourly" json:"priceHourly"`
	CloudProvider string  `db:"cloud_provider" json:"cloudProvider"`
	Vcpus         int     `db:"vcpus" json:"vcpus"`
	PriceMonthly  float64 `db:"price_monthly" json:"priceMonthly"`
	Region        string  `db:"region" json:"region"`
	Memory        int     `db:"memory" json:"memory"`
}

type NodesSlice []Nodes

func (n *Nodes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{n.NodeID, n.Description, n.Slug, n.Disk, n.PriceHourly, n.CloudProvider, n.Vcpus, n.PriceMonthly, n.Region, n.Memory}
	}
	return pgValues
}
func (n *Nodes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"node_id", "description", "slug", "disk", "price_hourly", "cloud_provider", "vcpus", "price_monthly", "region", "memory"}
	return columnValues
}
func (n *Nodes) GetTableName() (tableName string) {
	tableName = "nodes"
	return tableName
}

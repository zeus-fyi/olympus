package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Nodes struct {
	Memory        int     `db:"memory" json:"memory"`
	Vcpus         int     `db:"vcpus" json:"vcpus"`
	Disk          int     `db:"disk" json:"disk"`
	DiskUnits     string  `db:"disk_units" json:"diskUnits"`
	PriceHourly   float64 `db:"price_hourly" json:"priceHourly"`
	Region        string  `db:"region" json:"region"`
	CloudProvider string  `db:"cloud_provider" json:"cloudProvider"`
	ResourceID    int     `db:"resource_id" json:"resourceID"`
	Description   string  `db:"description" json:"description"`
	Slug          string  `db:"slug" json:"slug"`
	MemoryUnits   string  `db:"memory_units" json:"memoryUnits"`
	PriceMonthly  float64 `db:"price_monthly" json:"priceMonthly"`
}
type NodesSlice []Nodes

func (n *Nodes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{n.Memory, n.Vcpus, n.Disk, n.DiskUnits, n.PriceHourly, n.Region, n.CloudProvider, n.ResourceID, n.Description, n.Slug, n.MemoryUnits, n.PriceMonthly}
	}
	return pgValues
}
func (n *Nodes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"memory", "vcpus", "disk", "disk_units", "price_hourly", "region", "cloud_provider", "resource_id", "description", "slug", "memory_units", "price_monthly"}
	return columnValues
}
func (n *Nodes) GetTableName() (tableName string) {
	tableName = "nodes"
	return tableName
}

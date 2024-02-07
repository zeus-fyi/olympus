package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Nodes struct {
	ExtCfgStrID string `db:"ext_cfg_id" json:"extCfgStrID"`

	Memory        int     `db:"memory" json:"memory"`
	Vcpus         float64 `db:"vcpus" json:"vcpus"`
	Disk          int     `db:"disk" json:"disk"`
	DiskUnits     string  `db:"disk_units" json:"diskUnits"`
	DiskType      string  `db:"disk_type" json:"diskType"`
	PriceHourly   float64 `db:"price_hourly" json:"priceHourly"`
	Region        string  `db:"region" json:"region"`
	CloudProvider string  `db:"cloud_provider" json:"cloudProvider"`
	ResourceID    int     `db:"resource_id" json:"resourceID"`
	Description   string  `db:"description" json:"description"`
	Slug          string  `db:"slug" json:"slug"`
	MemoryUnits   string  `db:"memory_units" json:"memoryUnits"`
	PriceMonthly  float64 `db:"price_monthly" json:"priceMonthly"`
	Gpus          int     `db:"gpus" json:"gpus"`
	GpuType       string  `db:"gpu_type" json:"gpuType"`
}
type NodesSlice []Nodes

func (n *Nodes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{n.Memory, n.Vcpus, n.Disk, n.DiskUnits, n.DiskType, n.PriceHourly, n.Region, n.CloudProvider, n.ResourceID, n.Description, n.Slug, n.MemoryUnits, n.PriceMonthly, n.Gpus, n.GpuType}
	}
	return pgValues
}
func (n *Nodes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"memory", "vcpus", "disk", "disk_units", "disk_type", "price_hourly", "region", "cloud_provider", "resource_id", "description", "slug", "memory_units", "price_monthly", "gpus", "gpu_type"}
	return columnValues
}
func (n *Nodes) GetTableName() (tableName string) {
	tableName = "nodes"
	return tableName
}

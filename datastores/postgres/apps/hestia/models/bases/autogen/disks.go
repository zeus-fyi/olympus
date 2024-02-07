package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Disks struct {
	ExtCfgStrID string `db:"ext_cfg_id" json:"extCfgStrID"`

	ResourceID    int     `db:"resource_id" json:"resourceID"`
	DiskUnits     string  `db:"disk_units" json:"diskUnits"`
	PriceMonthly  float64 `db:"price_monthly" json:"priceMonthly"`
	Description   string  `db:"description" json:"description"`
	Type          string  `db:"type" json:"type"`
	SubType       string  `db:"sub_type" json:"subType"`
	DiskSize      int     `db:"disk_size" json:"diskSize"`
	PriceHourly   float64 `db:"price_hourly" json:"priceHourly"`
	Region        string  `db:"region" json:"region"`
	CloudProvider string  `db:"cloud_provider" json:"cloudProvider"`
}
type DisksSlice []Disks

func (d *Disks) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{d.ResourceID, d.DiskUnits, d.PriceMonthly, d.Description, d.Type, d.SubType, d.DiskSize, d.PriceHourly, d.Region, d.CloudProvider}
	}
	return pgValues
}
func (d *Disks) GetTableColumns() (columnValues []string) {
	columnValues = []string{"resource_id", "disk_units", "price_monthly", "description", "type", "sub_type", "disk_size", "price_hourly", "region", "cloud_provider"}
	return columnValues
}
func (d *Disks) GetTableName() (tableName string) {
	tableName = "disks"
	return tableName
}

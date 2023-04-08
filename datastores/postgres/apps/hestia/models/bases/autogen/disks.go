package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Disks struct {
	DiskID        int     `db:"disk_id" json:"diskID"`
	Description   string  `db:"description" json:"description"`
	Type          string  `db:"type" json:"type"`
	Size          int     `db:"size" json:"size"`
	PriceMonthly  float64 `db:"price_monthly" json:"priceMonthly"`
	Region        string  `db:"region" json:"region"`
	Units         string  `db:"units" json:"units"`
	PriceHourly   float64 `db:"price_hourly" json:"priceHourly"`
	CloudProvider string  `db:"cloud_provider" json:"cloudProvider"`
}
type DisksSlice []Disks

func (d *Disks) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{d.DiskID, d.Description, d.Type, d.Size, d.PriceMonthly, d.Region, d.Units, d.PriceHourly, d.CloudProvider}
	}
	return pgValues
}
func (d *Disks) GetTableColumns() (columnValues []string) {
	columnValues = []string{"disk_id", "description", "type", "size", "price_monthly", "region", "units", "price_hourly", "cloud_provider"}
	return columnValues
}
func (d *Disks) GetTableName() (tableName string) {
	tableName = "disks"
	return tableName
}

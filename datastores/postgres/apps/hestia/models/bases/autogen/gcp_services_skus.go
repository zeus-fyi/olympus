package hestia_autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type GcpServicesSkus struct {
	UsageType           sql.NullString `db:"usage_type" json:"usageType"`
	ServiceRegions      sql.NullString `db:"service_regions" json:"serviceRegions"`
	PricingInfo         sql.NullString `db:"pricing_info" json:"pricingInfo"`
	GeoTaxonomy         sql.NullString `db:"geo_taxonomy" json:"geoTaxonomy"`
	ServiceID           string         `db:"service_id" json:"serviceID"`
	ServiceDisplayName  sql.NullString `db:"service_display_name" json:"serviceDisplayName"`
	ResourceFamily      sql.NullString `db:"resource_family" json:"resourceFamily"`
	ResourceGroup       sql.NullString `db:"resource_group" json:"resourceGroup"`
	ServiceProviderName sql.NullString `db:"service_provider_name" json:"serviceProviderName"`
	Name                string         `db:"name" json:"name"`
	SkuID               string         `db:"sku_id" json:"skuID"`
	Description         sql.NullString `db:"description" json:"description"`
}
type GcpServicesSkusSlice []GcpServicesSkus

func (g *GcpServicesSkus) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{g.ServiceID, g.Name, g.SkuID, g.Description, g.ServiceDisplayName, g.ResourceFamily, g.ResourceGroup, g.UsageType, g.ServiceRegions, g.PricingInfo, g.ServiceProviderName, g.GeoTaxonomy}
	}
	return pgValues
}
func (g *GcpServicesSkus) GetTableColumns() (columnValues []string) {
	columnValues = []string{"usage_type", "service_regions", "pricing_info", "geo_taxonomy", "service_id", "service_display_name", "resource_family", "resource_group", "service_provider_name", "name", "sku_id", "description"}
	return columnValues
}
func (g *GcpServicesSkus) GetTableName() (tableName string) {
	tableName = "gcp_services_skus"
	return tableName
}

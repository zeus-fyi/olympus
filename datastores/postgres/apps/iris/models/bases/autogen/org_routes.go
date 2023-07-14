package iris_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgRoutes struct {
	RouteID   int    `db:"route_id" json:"routeID"`
	OrgID     int    `db:"org_id" json:"orgID"`
	RoutePath string `db:"route_path" json:"routePath"`
}
type OrgRoutesSlice []OrgRoutes

func (o *OrgRoutes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.RouteID, o.OrgID, o.RoutePath}
	}
	return pgValues
}
func (o *OrgRoutes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"route_id", "org_id", "route_path"}
	return columnValues
}
func (o *OrgRoutes) GetTableName() (tableName string) {
	tableName = "org_routes"
	return tableName
}

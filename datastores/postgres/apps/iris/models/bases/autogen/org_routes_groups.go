package iris_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgRoutesGroups struct {
	RouteGroupID int `db:"route_group_id" json:"routeGroupID"`
	RouteID      int `db:"route_id" json:"routeID"`
}
type OrgRoutesGroupsSlice []OrgRoutesGroups

func (o *OrgRoutesGroups) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.RouteGroupID, o.RouteID}
	}
	return pgValues
}
func (o *OrgRoutesGroups) GetTableColumns() (columnValues []string) {
	columnValues = []string{"route_group_id", "route_id"}
	return columnValues
}
func (o *OrgRoutesGroups) GetTableName() (tableName string) {
	tableName = "org_routes_groups"
	return tableName
}

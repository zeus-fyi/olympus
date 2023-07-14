package iris_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgRouteGroups struct {
	RouteGroupID   int    `db:"route_group_id" json:"routeGroupID"`
	OrgID          int    `db:"org_id" json:"orgID"`
	RouteGroupName string `db:"route_group_name" json:"routeGroupName"`
}
type OrgRouteGroupsSlice []OrgRouteGroups

func (o *OrgRouteGroups) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.RouteGroupID, o.OrgID, o.RouteGroupName}
	}
	return pgValues
}
func (o *OrgRouteGroups) GetTableColumns() (columnValues []string) {
	columnValues = []string{"route_group_id", "org_id", "route_group_name"}
	return columnValues
}
func (o *OrgRouteGroups) GetTableName() (tableName string) {
	tableName = "org_route_groups"
	return tableName
}

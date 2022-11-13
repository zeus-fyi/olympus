package autogen_bases

type OrgUsersTopologies struct {
	TopologyID int `db:"topology_id" json:"topologyID"`
	OrgID      int `db:"org_id" json:"orgID"`
	UserID     int `db:"user_id" json:"userID"`
}
type OrgUsersTopologiesSlice []OrgUsersTopologies

func (o *OrgUsersTopologies) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.TopologyID, o.OrgID, o.UserID}
	}
	return pgValues
}
func (o *OrgUsersTopologies) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "org_id", "user_id"}
	return columnValues
}
func (o *OrgUsersTopologies) GetTableName() (tableName string) {
	tableName = "org_users_topologies"
	return tableName
}

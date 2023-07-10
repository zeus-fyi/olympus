package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EventStatus struct {
	EventID int    `db:"event_id" json:"eventID"`
	Status  string `db:"status" json:"status"`
}
type EventStatusSlice []EventStatus

func (e *EventStatus) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.EventID, e.Status}
	}
	return pgValues
}
func (e *EventStatus) GetTableColumns() (columnValues []string) {
	columnValues = []string{"event_id", "status"}
	return columnValues
}
func (e *EventStatus) GetTableName() (tableName string) {
	tableName = "event_status"
	return tableName
}

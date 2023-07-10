package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Events struct {
	EventID int `db:"event_id" json:"eventID"`
}
type EventsSlice []Events

func (e *Events) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.EventID}
	}
	return pgValues
}
func (e *Events) GetTableColumns() (columnValues []string) {
	columnValues = []string{"event_id"}
	return columnValues
}
func (e *Events) GetTableName() (tableName string) {
	tableName = "events"
	return tableName
}

package artemis_autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type OrchestrationsScheduledToCloudCtxNs struct {
	OrchestrationScheduleID int       `db:"orchestration_schedule_id" json:"orchestrationScheduleID"`
	OrchestrationID         int       `db:"orchestration_id" json:"orchestrationID"`
	CloudCtxNsID            int       `db:"cloud_ctx_ns_id" json:"cloudCtxNsID"`
	Status                  string    `db:"status" json:"status"`
	DateScheduled           time.Time `db:"date_scheduled" json:"dateScheduled"`
}
type OrchestrationsScheduledToCloudCtxNsSlice []OrchestrationsScheduledToCloudCtxNs

func (o *OrchestrationsScheduledToCloudCtxNs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrchestrationScheduleID, o.OrchestrationID, o.CloudCtxNsID, o.DateScheduled}
	}
	return pgValues
}
func (o *OrchestrationsScheduledToCloudCtxNs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"orchestration_schedule_id", "orchestration_id", "cloud_ctx_ns_id", "date_scheduled"}
	return columnValues
}
func (o *OrchestrationsScheduledToCloudCtxNs) GetTableName() (tableName string) {
	tableName = "orchestrations_scheduled_to_cloud_ctx_ns"
	return tableName
}

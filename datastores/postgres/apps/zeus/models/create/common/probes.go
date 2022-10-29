package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers/probes"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateProbeValueSubCTEs(containerID int, probes probes.ProbeSlice) (sql_query_templates.SubCTEs, sql_query_templates.SubCTEs) {
	var ts chronos.Chronos

	probesValuesSubCTEs := make([]sql_query_templates.SubCTE, len(probes))
	probesValuesRelationshipSubCTEs := make([]sql_query_templates.SubCTE, len(probes))
	for i, pr := range probes {
		timeNow := ts.UnixTimeStampNow()
		pr.SetProbeID(timeNow)
		probesValuesSubCTEs[i] = createProbeValueSubCTEs(&pr, i)
		probesValuesRelationshipSubCTEs[i] = createProbeRelationshipSubCTEs(containerID, &pr, i)
	}
	return probesValuesSubCTEs, probesValuesRelationshipSubCTEs
}

func createProbeValueSubCTEs(probe *probes.Probe, count int) sql_query_templates.SubCTE {
	var ts chronos.Chronos
	queryName := fmt.Sprintf("cte_%s_%d_%d_value", probe.ContainersProbes.GetTableName(), ts.UnixTimeStampNow(), count)
	subCTE := sql_query_templates.NewSubInsertCTE(queryName)
	subCTE.TableName = probe.ContainerProbes.GetTableName()
	subCTE.Columns = probe.ContainerProbes.GetTableColumns()
	subCTE.Values = []apps.RowValues{probe.ContainerProbes.GetRowValues(queryName)}
	return subCTE
}

func createProbeRelationshipSubCTEs(containerID int, probe *probes.Probe, count int) sql_query_templates.SubCTE {
	var ts chronos.Chronos
	queryName := fmt.Sprintf("cte_%s_%d_%d_value", probe.ContainersProbes.GetTableName(), ts.UnixTimeStampNow(), count)
	probe.SetContainerID(containerID)
	subCTE := sql_query_templates.NewSubInsertCTE(queryName)
	subCTE.TableName = probe.ContainersProbes.GetTableName()
	subCTE.Columns = probe.ContainersProbes.GetTableColumns()
	subCTE.Values = []apps.RowValues{probe.ContainersProbes.GetRowValues(queryName)}
	return subCTE
}

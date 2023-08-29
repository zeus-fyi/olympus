package iris_redis

import (
	"fmt"
	"time"

	util "github.com/wealdtech/go-eth2-util"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func getOrgRouteKey(orgID int, rgName string) string {
	return fmt.Sprintf("{%d}.%s", orgID, rgName)
}

func getOrgMonthlyUsageKey(orgID int, month string) string {
	return fmt.Sprintf("{%d}.%s-total-zu-usage", orgID, month)
}

func getOrgTotalRequestsKey(orgID int) string {
	return fmt.Sprintf("{%d}.total-zu-usage", orgID)
}

func getOrgRateLimitKey(orgID int) string {
	return fmt.Sprintf("{%d}.%d", orgID, time.Now().Unix())
}

func getTableMetricKey(orgID int, tableName, metric string) string {
	return fmt.Sprintf("{%d}:%s:%s", orgID, tableName, metric)
}

func getTableMetricSetKey(orgID int, tableName string) string {
	return fmt.Sprintf("{%d}:%s:metrics", orgID, tableName)
}

func getMetricTdigestKey(orgID int, tableName, metricName string) string {
	return fmt.Sprintf("{%d}.%s:%s", orgID, tableName, metricName)
}

func getMetricTdigestMetricSamplesKey(orgID int, tableName, metricName string) string {
	return fmt.Sprintf("%s:%s", getMetricTdigestKey(orgID, tableName, metricName), "samples")
}

func createAdaptiveEndpointPriorityScoreKey(orgID int, tableName string) string {
	return fmt.Sprintf("{%d}:%s:priority", orgID, tableName)
}

func getHashedTokenKey(token string) string {
	return fmt.Sprintf("{%x}", util.Keccak256([]byte(token)))
}

func getHashedTokenPlanKey(token string) string {
	return fmt.Sprintf("{%x}:plan", util.Keccak256([]byte(token)))
}

func getProcedureKey(orgID int, procedureName string) string {
	if orgID > 0 && procedureName != iris_programmable_proxy_v1_beta.MaxBlockAggReduce {
		return fmt.Sprintf("{%d}:%s:procedure", orgID, procedureName)
	}
	return getGlobalProcedureKey(procedureName)
}

func getProcedureStepsKey(orgID int, procedureName string) string {
	if orgID > 0 && procedureName != iris_programmable_proxy_v1_beta.MaxBlockAggReduce {
		return fmt.Sprintf("%s:steps", getProcedureKey(orgID, procedureName))
	}
	return getGlobalProcedureStepsKey(procedureName)
}

func getGlobalProcedureKey(procedureName string) string {
	return fmt.Sprintf("{global}.%s:procedure", procedureName)
}

func getGlobalProcedureStepsKey(procedureName string) string {
	return fmt.Sprintf("%s:steps", getGlobalProcedureKey(procedureName))
}

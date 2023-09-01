import {TableMetric} from "../../../../redux/loadbalancing/loadbalancing.types";

export const getSummaryStatsExt = (data: TableMetric) => {
    let min = 0
    let q1 = 0
    let median = 0
    let q3 = 0
    let max = 0
    for (const [key, tableMetrics] of Object.entries(data.metricPercentiles)) {

            if (tableMetrics.percentile == 0.1) {
                 min =tableMetrics.latency
            }
            if (tableMetrics.percentile == 0.25) {
                 q1 = tableMetrics.latency
            }
            if (tableMetrics.percentile == 0.5) {
                 median = tableMetrics.latency
            }
            if (tableMetrics.percentile == 0.75) {
                 q3 = tableMetrics.latency
            }
            if (tableMetrics.percentile == 0.99) {
                 max = tableMetrics.latency
            }
    }

    if(!q3 || !q1 || !median){
        return
    }

    const interQuantileRange = q3 - q1
    const minAdj = q1 - 1.5 * interQuantileRange
    const maxAdj = q3 + 1.5 * interQuantileRange

    return {minAdj, q1, median, q3, maxAdj}
}

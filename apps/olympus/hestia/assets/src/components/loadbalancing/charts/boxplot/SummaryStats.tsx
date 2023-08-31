import {TableMetricsSummary} from "../../../../redux/loadbalancing/loadbalancing.types";

export const getSummaryStatsExt = (data: TableMetricsSummary) => {
    let min = 0
    let q1 = 0
    let median = 0
    let q3 = 0
    let max = 0
    for (const [key, tableMetrics] of Object.entries(data.metrics)) {
        if (tableMetrics.metricPercentiles.length === 0) {
            return
        }
        if (tableMetrics.sampleCount === 0) {
            return
        }
        for (const [mpk, mp] of Object.entries(tableMetrics.metricPercentiles)) {
            if (mp.percentile == 0.1) {
                 min = mp.latency
            }
            if (mp.percentile == 0.25) {
                 q1 = mp.latency
            }
            if (mp.percentile == 0.5) {
                 median = mp.latency
            }
            if (mp.percentile == 0.75) {
                 q3 = mp.latency
            }
            if (mp.percentile == 0.99) {
                 max = mp.latency
            }
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

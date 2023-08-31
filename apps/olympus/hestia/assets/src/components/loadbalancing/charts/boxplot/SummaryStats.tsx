import * as d3 from "d3";
import {TableMetricsSummary} from "../../../../redux/loadbalancing/loadbalancing.types";

// Takes an array of numbers and compute some summary statistics from it like quantiles, median..
// Those summary statistics are the info needed to draw a boxplot
export const getSummaryStats = (data: number[]) => {
    const sortedData = data.sort(function(a, b){return a - b});

    const q1 = d3.quantile(sortedData, .25)
    const median = d3.quantile(sortedData, .5)
    const q3 = d3.quantile(sortedData, .75)

    if(!q3 || !q1 || !median){
        return
    }
    const interQuantileRange = q3 - q1
    const min = q1 - 1.5 * interQuantileRange
    const max = q3 + 1.5 * interQuantileRange

    return {min, q1, median, q3, max}
}

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

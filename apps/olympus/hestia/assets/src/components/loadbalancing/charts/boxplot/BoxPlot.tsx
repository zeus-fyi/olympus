import * as d3 from "d3";
import {AxisLeft} from "./AxisLeft";
import {AxisBottom} from "./AxisBottomCategoric";
import {TableMetricsSummary} from "../../../../redux/loadbalancing/loadbalancing.types";
import {getSummaryStatsExt} from "./SummaryStats";
import {VerticalBox} from "./VerticalBox";

const MARGIN = { top: 175, right: 175, bottom: 30, left: 75 };

type BoxplotProps = {
    width: number;
    height: number;
    tableMetrics: TableMetricsSummary;
};

export const Boxplot = ({ width, height, tableMetrics }: BoxplotProps) => {
    if (!tableMetrics) {
        return null;
    }

    // Scales and Dimensions
    const boundsWidth = width - MARGIN.right - MARGIN.left;
    const boundsHeight = height - MARGIN.top - MARGIN.bottom;

    // Data preprocessing
    const groups = Object.keys(tableMetrics.metrics);
    const allLatencies = Object.values(tableMetrics.metrics).flatMap((metric) =>
        metric.metricPercentiles.map((sample) => sample.latency)
    );
    const maxVal = d3.max(allLatencies) ?? 0;
    const minVal = d3.min(allLatencies) ?? 0;
    // Define scales
    const yScale = d3.scaleBand().range([0, boundsHeight]).domain(groups);
    const xScale = d3.scaleLinear().range([0, boundsWidth]).domain([minVal*0.75, maxVal*1.25]);

    // Render BoxPlots
    const allShapes = Object.entries(tableMetrics.metrics).map(([key, metric], i) => {
        const sumStats = getSummaryStatsExt(metric);  // Get summary stats for each metric
        if (!sumStats) {
            return null;
        }
        const { minAdj, q1, median, q3, maxAdj,max,min } = sumStats;
        return (
            <g key={key} transform={`translate(0, ${yScale(key) ?? 0})`}>
                <VerticalBox
                    width={q3-q1}
                    min={xScale(min) ?? 0}
                    q1={xScale(q1) ?? 0}
                    median={xScale(median) ?? 0}
                    q3={xScale(q3) ?? 0}
                    max={xScale(max) ?? 0}
                    stroke="black"
                    fill="#ead4f5"
                    offset={i}
                />
            </g>
        );
    });

    const pixelPerTick = boundsHeight / groups.length+1;
    return (
        <div>
            <svg width={width} height={height}>
                <g transform={`translate(${MARGIN.left+150}, ${MARGIN.top})`}>
                    {allShapes}
                    <AxisLeft yScale={yScale} pixelsPerTick={pixelPerTick} height={boundsHeight} />
                    <g transform={`translate(0, ${boundsHeight})`}>
                        <AxisBottom xScale={xScale} />
                    </g>
                </g>
            </svg>
        </div>
    );
};
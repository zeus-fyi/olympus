import * as d3 from "d3";
import {AxisLeft} from "./AxisLeft";
import {AxisBottom} from "./AxisBottomCategoric";
import {TableMetricsSummary} from "../../../../redux/loadbalancing/loadbalancing.types";
import {getSummaryStatsExt} from "./SummaryStats";
import {VerticalBox} from "./VerticalBox";

const MARGIN = { top: 30, right: 30, bottom: 30, left: 50 };

type BoxplotProps = {
    width: number;
    height: number;
    data: { name: string; value: number }[];
    tableMetrics: TableMetricsSummary;
};

export const Boxplot = ({ width, height, data, tableMetrics }: BoxplotProps) => {
    if (!tableMetrics) {
        return null;
    }

    const sumStats = getSummaryStatsExt(tableMetrics);
    if (!sumStats) {
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
    const maxVal = d3.max(allLatencies);

    // Define scales
    const yScale = d3.scaleBand().range([0, boundsHeight]).domain(groups).padding(0.1); // Added padding
    const xScale = d3.scaleLinear().range([0, boundsWidth]).domain([0, maxVal ?? 0]);

    // Render BoxPlots
    const allShapes = Object.entries(tableMetrics.metrics).map(([key, metric], i) => {
        const { minAdj, q1, median, q3, maxAdj } = sumStats;

        return (
            <g key={key} transform={`translate(0, ${yScale(key) ?? 0})`}>
                <VerticalBox
                    width={boundsWidth}
                    min={xScale(minAdj) ?? 0}
                    q1={xScale(q1) ?? 0}
                    median={xScale(median) ?? 0}
                    q3={xScale(q3) ?? 0}
                    max={xScale(maxAdj) ?? 0}
                    stroke="black"
                    fill="#ead4f5"
                />
            </g>
        );
    });

    return (
        <div>
            <svg width={width} height={height}>
                <g transform={`translate(${MARGIN.left}, ${MARGIN.top})`}>
                    {allShapes}
                    <AxisLeft yScale={yScale} pixelsPerTick={30} />
                    <g transform={`translate(0, ${boundsHeight})`}>
                        <AxisBottom xScale={xScale} />
                    </g>
                </g>
            </svg>
        </div>
    );
};

import {useMemo} from "react";
import * as d3 from "d3";
import {AxisLeft} from "./AxisLeft";
import {AxisBottom} from "./AxisBottomCategoric";
import {TableMetric, TableMetricsSummary} from "../../../../redux/loadbalancing/loadbalancing.types";
import {getSummaryStatsExt} from "./SummaryStats";
import {VerticalBox} from "./VerticalBox";

const MARGIN = { top: 30, right: 30, bottom: 30, left: 50 };

type BoxplotProps = {
    width: number;
    height: number;
    data: { name: string; value: number }[];
    tableMetrics : TableMetricsSummary;
};

export const Boxplot = ({ width, height, data, tableMetrics }: BoxplotProps) => {
    // Compute everything derived from the dataset:
    const { groups } = useMemo(() => {
        // @ts-ignore
        const groups = Object.entries(tableMetrics.metrics).flatMap(([key, metric]: [string, TableMetric]) => {
            // Extract latencies from the metricPercentiles array
            return key
        });
        return { groups };
    }, [data]);
    if (tableMetrics == undefined ) {
        return null;
    }
    const sumStats = getSummaryStatsExt(tableMetrics);
    if (!sumStats) {
        return null;
    }
    console.log("Boxplot.tsx: Boxplot: data: ", tableMetrics);
    // The bounds (= area inside the axis) is calculated by substracting the margins from total width / height
    const { minAdj, q1, median, q3, maxAdj } = sumStats;
    const boundsWidth = width - MARGIN.right - MARGIN.left;
    const boundsHeight = height - MARGIN.top - MARGIN.bottom;
    // Compute scales
    const tm = Object.entries(tableMetrics.metrics).flatMap(([key, metric]: [string, TableMetric]) => {
        console.log("Boxplot.tsx: Boxplot: key: ", key);
        console.log("Boxplot.tsx: Boxplot: metric: ", metric);
        // Extract latencies from the metricPercentiles array
        return metric.metricPercentiles.map(sample => sample.latency);
    });

    const filteredTm = tm.filter((t): t is number => t !== undefined);
    const maxVal = d3.max(filteredTm);

    console.log("Boxplot.tsx: Boxplot: tm: ", tm);
    const yScale = d3
        .scaleBand()
        .range([boundsHeight, 0])
        .domain(groups)
        .padding(0);
    const xScale = d3
        .scaleLinear()
        .range([0, boundsWidth])
        .domain([0, maxVal ?? 1])

    // domain should be be the range of numbers in a slice passed to xScale I guess
    // Build the box shapes
    const allShapes = groups.map((group, i) => {
        return (
            <g key={i} transform={`translate(${yScale(group)},0)`}>
                <VerticalBox
                    width={yScale.bandwidth()}
                    q1={xScale(q1)}
                    median={xScale(median)}
                    q3={xScale(q3)}
                    min={xScale(minAdj)}
                    max={xScale(maxAdj)}
                    stroke="black"
                    fill={"#ead4f5"}
                />
            </g>
        );
    });

    return (
        <div>
            <svg width={width} height={height}>
                <g
                    width={boundsWidth}
                    height={boundsHeight}
                    transform={`translate(${[MARGIN.left, MARGIN.top].join(",")})`}
                >
                    {/*{allShapes}*/}
                    <AxisLeft yScale={yScale} pixelsPerTick={30} />
                    {/* X axis uses an additional translation to appear at the bottom */}
                    <g transform={`translate(0, ${boundsHeight})`}>
                        <AxisBottom xScale={xScale} />
                    </g>
                </g>
            </svg>
        </div>
    );
};

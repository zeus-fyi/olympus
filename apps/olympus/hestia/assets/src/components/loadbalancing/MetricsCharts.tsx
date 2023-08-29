import React from 'react';
import {PolarAngleAxis, PolarGrid, PolarRadiusAxis, Radar, RadarChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {TableMetric, TableMetricsSummary} from "../../redux/loadbalancing/loadbalancing.types";


interface MetricsChartProps {
    tableMetricsSummary: TableMetricsSummary;
}

interface FormattedData {
    subject: string;
    sampleCount: number;
    percentile: number;
    latency: number;
}

export function TableMetricsCharts(props: any) {
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    const formattedData: FormattedData[] = Object.keys(tableMetrics.metrics).map((key, idx) => {
        const metric: TableMetric = tableMetrics.metrics[key];
        return {
            subject: `Sample ${idx + 1}`,
            sampleCount: metric.sampleCount,
            percentile: metric.metricPercentiles.reduce((acc, sample) => acc + sample.percentile, 0) / metric.metricPercentiles.length,
            latency: metric.metricPercentiles.reduce((acc, sample) => acc + sample.latency, 0) / metric.metricPercentiles.length,
        };
    });
    console.log('formattedData', formattedData)

    return (
        <div>
            <MetricsChart formattedData={formattedData}/>
        </div>
    )
}

export function MetricsChart(props: any) {
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    const formattedData: FormattedData[] = Object.keys(tableMetrics.metrics).map((key, idx) => {
        const metric: TableMetric = tableMetrics.metrics[key];
        return {
            subject: `Sample ${idx + 1}`,
            sampleCount: metric.sampleCount,
            percentile: metric.metricPercentiles.reduce((acc, sample) => acc + sample.percentile, 0) / metric.metricPercentiles.length,
            latency: metric.metricPercentiles.reduce((acc, sample) => acc + sample.latency, 0) / metric.metricPercentiles.length,
        };
    });
    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                    Metrics Chart for {tableMetrics.tableName}
                </Typography>
                <RadarChart cx="50%" cy="50%" outerRadius="80%" data={formattedData}>
                    <PolarGrid />
                    <PolarAngleAxis dataKey="subject" />
                    <PolarRadiusAxis />
                    <Radar name="Sample Count" dataKey="sampleCount" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
                    <Radar name="Percentile" dataKey="percentile" stroke="#82ca9d" fill="#82ca9d" fillOpacity={0.6} />
                    <Radar name="Latency" dataKey="latency" stroke="#ffc658" fill="#ffc658" fillOpacity={0.6} />
                </RadarChart>
            </CardContent>
        </Card>
    )
}

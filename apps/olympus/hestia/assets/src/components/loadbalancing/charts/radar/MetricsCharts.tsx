import React, {useEffect} from 'react';
import {PolarAngleAxis, PolarGrid, PolarRadiusAxis, Radar, RadarChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import Container from "@mui/material/Container";
import {RootState} from "../../../../redux/store";
import {Boxplot} from "../boxplot/BoxPlot";
import {TableMetric} from "../../../../redux/loadbalancing/loadbalancing.types";
import {loadBalancingApiGateway} from "../../../../gateway/loadbalancing";
import {setTableMetrics} from "../../../../redux/loadbalancing/loadbalancing.reducer";

interface FormattedData {
    subject: string;
    sampleCount: number;
    percentile: number;
    latency: number;
}

export function TableMetricsCharts(props: any) {
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    const dispatch = useDispatch();
    const [loading, setLoading] = React.useState(false);
    const {tableName} = props;

    useEffect(() => {
        async function fetchData() {
            try {
                setLoading(true);
                const response = await loadBalancingApiGateway.getTableMetrics(tableName);
                console.log(response.data)
                console.log(response)
                const tableMetrics = response.data;
                if (tableMetrics.metrics != null && tableMetrics.length > 0) {
                    dispatch(setTableMetrics(response.data));
                }
            } catch (e) {
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, [tableName]);

    if (loading) {
        return <div>Loading...</div>
    }

    return (
        <div>
            {/*<MetricsChart />*/}
            {tableMetrics && tableMetrics.metrics &&
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Boxplot tableMetrics={tableMetrics} width={1200} height={800} />
            </Container>}
        </div>
    )
}

export function MetricsChart(props: any) {
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    if (tableMetrics == null || tableMetrics.metrics == null ||  Object.keys(tableMetrics.metrics).length == 0) {
        return <div></div>
    }
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
        <div>
            {formattedData.length > 0 && (
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
            )}
        </div>
    )
}

import React, {useEffect} from 'react';
import {PolarAngleAxis, PolarGrid, Radar, RadarChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import Container from "@mui/material/Container";
import {RootState} from "../../../../redux/store";
import {Boxplot} from "../boxplot/BoxPlot";
import {generateMetricSlices, MetricAggregateRow} from "../../../../redux/loadbalancing/loadbalancing.types";
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
                // console.log(response.data)
                // console.log(response)
                const tableMetrics = response.data;
                if (tableMetrics != null && tableMetrics.metrics != null) {
                    dispatch(setTableMetrics(tableMetrics));
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

    if (tableMetrics == null || tableMetrics.metrics == null ||  Object.keys(tableMetrics.metrics).length == 0) {
        return <div></div>
    }
    const metricSlices: MetricAggregateRow[] = generateMetricSlices([tableMetrics]); // Generate slices here
    let safeEndpoints = metricSlices ?? [];
    if (safeEndpoints.length == 0) {
        return <div></div>
    }
    return (
        <div>
            {/*<MetricsChart />*/}
            {tableMetrics && tableMetrics.metrics &&
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Boxplot tableMetrics={tableMetrics} width={1200} height={(1+safeEndpoints.length)*200} />
            </Container>}
        </div>
    )
}

export function MetricsChart(props: any) {
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    if (tableMetrics == null || tableMetrics.metrics == null) {
        return <div></div>
    }

    const formattedData: MetricAggregateRow[] = generateMetricSlices([tableMetrics]); // Generate slices here
    return (
        <div>
            {formattedData.length > 0 && (
                <Card>
                    <CardContent>
                        <Typography variant="h5" gutterBottom>
                            Table Requests
                        </Typography>
                        <RadarChart cx="50%" cy="50%" outerRadius="80%" width={700} height={300} data={formattedData}>
                            <PolarGrid />
                            <PolarAngleAxis dataKey="metricName" />
                            <Radar name="Sample Count" dataKey="sampleCount" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
                        </RadarChart>
                    </CardContent>
                </Card>
            )}
        </div>
    )
}

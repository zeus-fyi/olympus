import React, {useState} from 'react';
import {Legend, Pie, PieChart} from 'recharts';
import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";

export function PlanUsagePieCharts(props: any) {
    const planUsageDetails = useSelector((state: RootState) => state.loadBalancing.planUsageDetails);
    return (
        <div>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                <PlanRateUsagePieChart  planUsageDetails={planUsageDetails}/>
                <PlanTableCountUsagePieChart planUsageDetails={planUsageDetails} />
            </Stack>
        </div>
    )
}

export function PlanRateUsagePieChart(props: any) {
    const { planUsageDetails, reload, setReload} = props;
    const title = planUsageDetails?.planName +  ' Plan';
    const [loading, setLoading] = useState(true);
    const [rateLimit, setRateLimit] = useState(planUsageDetails?.computeUsage?.rateLimit ?? 0);
    const [currentRate, setCurrentRate] = useState(planUsageDetails?.computeUsage?.currentRate ?? 0);
    const [planBudgetZU, setPlanBudgetZU] = useState(planUsageDetails?.computeUsage?.monthlyBudgetZU ?? 0);
    const [monthlyUsage, setMonthlyUsage] = useState(planUsageDetails?.computeUsage?.monthlyUsage ?? 0);
    const data02 = [
        { name: 'ZU k/s', value: currentRate, fill: "#ff8080"},
        { name: 'ZU k/s limit', value: rateLimit - currentRate, fill: "#82ca9d"},
    ];
    const data01 = [
        { name: 'ZU M remaining', value: planBudgetZU-monthlyUsage, fill: "#4db375"},
        { name: 'ZU M consumed', value: monthlyUsage, fill: "#ff4d4d"},
    ];

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                    {title}
                </Typography>
                <PieChart width={375} height={275}>
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60}  />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} label />
                    <Legend align="left" verticalAlign="bottom" layout="horizontal" />
                </PieChart>
            </CardContent>
        </Card>
    )
}


export function PlanTableCountUsagePieChart(props: any) {
    const {planUsageDetails, reload, setReload} = props;
    const dispatch = useDispatch();
    const [loading, setLoading] = useState(true);
    const [endpointCount, setEndpointCount] = useState(planUsageDetails?.tableUsage?.endpointCount);
    const maxEndpointCount = 1000;
    const remainingEndpoints = maxEndpointCount - endpointCount;
    const [tableCount, setTableCount] = useState(planUsageDetails?.tableUsage?.tableCount);
    const [planTableCount, setPlanTableCount] = useState(planUsageDetails?.tableUsage?.monthlyBudgetTableCount);
    const remainingTables = planTableCount - tableCount;
    const data01 = [
        { name: 'Endpoints(Used)', value: endpointCount, fill: "#8884d8" },
        { name: 'Endpoints(Open)', value: remainingEndpoints, fill: "#82ca9d" },
    ];
    const data02 = [
        { name: 'Tables(Used)', value: endpointCount, fill: "#8884d8" },
        { name: 'Tables(Open)', value: remainingTables, fill: "#82ca9d" },
    ];

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                   Table, Endpoint Usage
                </Typography>
                <PieChart width={375} height={275}>
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60} fill="#8884d8" />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} fill="#82ca9d" label />
                    <Legend align="left" verticalAlign="bottom" layout="horizontal" />
                </PieChart>
            </CardContent>
        </Card>
    )
}
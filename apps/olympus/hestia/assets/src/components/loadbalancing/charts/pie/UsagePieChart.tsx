import React, {useState} from 'react';
import {Legend, Pie, PieChart} from 'recharts';
import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";

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
    const [rateLimit, setRateLimit] = useState(planUsageDetails?.computeUsage?.rateLimit ?? 0);
    const [currentRate, setCurrentRate] = useState(planUsageDetails?.computeUsage?.currentRate ?? 0);
    const [planBudgetZU, setPlanBudgetZU] = useState(planUsageDetails?.computeUsage?.monthlyBudgetZU ?? 0);
    const [monthlyUsage, setMonthlyUsage] = useState(planUsageDetails?.computeUsage?.monthlyUsage ?? 0);

    const data02 = [
        { name: `ZU ${(rateLimit - currentRate).toFixed(2)}k/s  limit`, value: (rateLimit - currentRate).toFixed(2), fill: "#4db375"},
        { name: `ZU ${(currentRate).toFixed(2)} k/s`, value: currentRate.toFixed(2), fill: "#ff4d4d"},
    ];
    const data01 = [
        { name: `ZU ${(planBudgetZU-monthlyUsage).toFixed(2)}M`, value: (planBudgetZU-monthlyUsage), fill: "#4db375"},
        { name: `ZU ${(monthlyUsage).toFixed(2)}M used`, value: monthlyUsage, fill: "#ff4d4d"},
    ];

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                    {title}
                </Typography>
                <PieChart width={450} height={275}>
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60}  />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} label />
                    <Legend align="left" verticalAlign="bottom" layout="vertical" />
                </PieChart>
            </CardContent>
        </Card>
    )
}


export function PlanTableCountUsagePieChart(props: any) {
    const {planUsageDetails, reload, setReload} = props;
    const [endpointCount, setEndpointCount] = useState(planUsageDetails?.tableUsage?.endpointCount);
    const maxEndpointCount = 1000;
    const remainingEndpoints = maxEndpointCount - endpointCount;
    const [tableCount, setTableCount] = useState(planUsageDetails?.tableUsage?.tableCount);
    const [planTableCount, setPlanTableCount] = useState(planUsageDetails?.tableUsage?.monthlyBudgetTableCount);
    const remainingTables = planTableCount - tableCount;

    const data01 = [
        { name: `routes open ${remainingEndpoints.toFixed(0)}`, value: remainingEndpoints, fill: "#4db375" },
        { name: `routes used ${endpointCount.toFixed(0)}`, value: endpointCount, fill: "#ff4d4d" },
    ];
    const data02 = [
        { name: `tables open ${remainingTables.toFixed(0)}`, value: remainingTables, fill: "#4db375" },
        { name: `tables used ${tableCount.toFixed(0)}`, value: tableCount, fill: "#ff4d4d" },
    ];

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                   Table & Route Usage
                </Typography>
                <PieChart width={450} height={275}>
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60} fill="#8884d8" />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} fill="#82ca9d" label />
                    <Legend align="left" verticalAlign="bottom" layout="vertical" />
                </PieChart>
            </CardContent>
        </Card>
    )
}
import React from 'react';
import {Legend, Pie, PieChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";

const data01 = [
    { name: 'Usage A', value: 100 },
    { name: 'Usage B', value: 300 },
];
const data02 = [
    { name: 'A1', value: 100 },
    { name: 'A2', value: 300 },
];


export function PlanUsagePieChart(props: any) {
    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                    Rate Usage
                </Typography>
                <PieChart width={375} height={275}>
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60} fill="#8884d8" />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} fill="#82ca9d" label />
                    <Legend />
                </PieChart>
            </CardContent>
        </Card>
    )
}
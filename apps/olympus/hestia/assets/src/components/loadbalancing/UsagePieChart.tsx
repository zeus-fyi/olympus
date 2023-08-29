import React, {useEffect, useState} from 'react';
import {Legend, Pie, PieChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import {accessApiGateway} from "../../gateway/access";
import {setSessionAuth} from "../../redux/auth/session.reducer";
import {setUserPlanDetails} from "../../redux/loadbalancing/loadbalancing.reducer";
import {useDispatch} from "react-redux";

const data01 = [
    { name: 'Consumed', value: 100 },
    { name: 'Remaining', value: 300 },
];
const data02 = [
    { name: 'Endpoints', value: 100 },
    { name: 'Unused', value: 300 },
];

export function PlanUsagePieChart(props: any) {
    const {planUsageDetails} = props;
    const title = planUsageDetails?.planName +  ' Plan';
    const dispatch = useDispatch();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                if (planUsageDetails?.planName !== null && planUsageDetails?.planName !== ''){
                        return
                }
                const response = await accessApiGateway.checkAuth();
                if (response.status !== 200) {
                    dispatch(setSessionAuth(false));
                    return;
                }
                if (response.data.planUsageDetails != null){
                    dispatch(setUserPlanDetails(response.data.planUsageDetails))
                }
            } catch (error) {
                setLoading(false);
            } finally {
                setLoading(false);
            }
        }
        fetchData().then(r =>
            console.log("")
        );
    }, []);

    if (loading) {
        return null;
    }
    return (
        <Card>
            <CardContent>
                <Typography variant="h5" gutterBottom>
                    {title}
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
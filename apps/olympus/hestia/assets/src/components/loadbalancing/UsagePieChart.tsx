import React, {useEffect, useState} from 'react';
import {Legend, Pie, PieChart} from 'recharts';
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import {accessApiGateway} from "../../gateway/access";
import {setSessionAuth} from "../../redux/auth/session.reducer";
import {setUserPlanDetails} from "../../redux/loadbalancing/loadbalancing.reducer";
import {useDispatch} from "react-redux";


export function PlanRateUsagePieChart(props: any) {
    const {planUsageDetails} = props;
    const title = planUsageDetails?.planName +  ' Plan';
    const dispatch = useDispatch();
    const [loading, setLoading] = useState(true);
    const [rateLimit, setRateLimit] = useState(planUsageDetails?.computeUsage?.rateLimit);
    const [currentRate, setCurrentRate] = useState(planUsageDetails?.computeUsage?.currentRate);
    const remainingRate = rateLimit - currentRate;
    const [planBudgetZU, setPlanBudgetZU] = useState(planUsageDetails?.computeUsage?.monthlyBudgetZU);
    const [monthlyUsage, setMonthlyUsage] = useState(planUsageDetails?.computeUsage?.monthlyUsage);
    const remainingZU = planBudgetZU - monthlyUsage;
    const data02 = [
        { name: 'ZU k/s', value: currentRate, fill: "#82ca9d"},
        { name: 'ZU k/s limit', value: remainingRate, fill: "#ff8080"},
    ];
    const data01 = [
        { name: 'ZU M remaining', value: remainingZU, fill: "#4db375"},
        { name: 'ZU M consumed', value: monthlyUsage, fill: "#ff4d4d"},
    ];

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
                    <Pie data={data01} dataKey="value" cx="50%" cy="50%" outerRadius={60}  />
                    <Pie data={data02} dataKey="value" cx="50%" cy="50%" innerRadius={70} outerRadius={90} label />
                    <Legend />
                </PieChart>
            </CardContent>
        </Card>
    )
}


export function PlanTableCountUsagePieChart(props: any) {
    const {planUsageDetails} = props;
    const dispatch = useDispatch();
    const [loading, setLoading] = useState(true);
    const [endpointCount, setEndpointCount] = useState(planUsageDetails?.tableUsage?.endpointCount);
    const maxEndpointCount = 1000;
    const remainingEndpoints = maxEndpointCount - endpointCount;
    const [tableCount, setTableCount] = useState(planUsageDetails?.tableUsage?.tableCount);
    const [planTableCount, setPlanTableCount] = useState(planUsageDetails?.tableUsage?.monthlyBudgetTableCount);
    const remainingTables = planTableCount - tableCount;
    const data01 = [
        { name: 'Tables(Used)', value: endpointCount },
        { name: 'Tables(Open)', value: remainingTables },
    ];
    const data02 = [
        { name: 'Endpoints(Used)', value: endpointCount },
        { name: 'Endpoints(Open)', value: remainingEndpoints },
    ];
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
                   Table, Endpoint Usage
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
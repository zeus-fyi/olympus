import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Groups, LoadBalancingState, PlanUsageDetails, TableMetricsSummary} from "./loadbalancing.types";

const initialState: LoadBalancingState = {
    routes: [],
    groups: {},
    planUsageDetails: {
        planName: '',
        computeUsage: {
            rateLimit: 0,
            currentRate: 0,
            monthlyUsage: 0,
            monthlyBudgetZU: 0
        },
        tableUsage: {
            monthlyBudgetTableCount: 0,
            endpointCount: 0,
            tableCount: 0
        }
    },
    tableMetrics: {
        tableName: '',
        routes: [],
        metrics: {},
    },
}

const loadBalancingSlice = createSlice({
    name: 'loadBalancing',
    initialState,
    reducers: {
        setEndpoints: (state, action: PayloadAction<any[]>) => {
            state.routes = action.payload;
        },
        setGroupEndpoints: (state, action: PayloadAction<Groups>) => {
            state.groups = action.payload;
        },
        setUserPlanDetails: (state, action: PayloadAction<PlanUsageDetails>) => {
            state.planUsageDetails = action.payload;
        },
        setTableMetrics: (state, action: PayloadAction<TableMetricsSummary>) => {
            state.tableMetrics = action.payload;
        }
    }
});

export const {
    setEndpoints,
    setGroupEndpoints,
    setUserPlanDetails,
    setTableMetrics
} = loadBalancingSlice.actions;
export default loadBalancingSlice.reducer;
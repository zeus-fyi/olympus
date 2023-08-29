import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Groups, LoadBalancingState, PlanUsageDetails} from "./loadbalancing.types";

const initialState: LoadBalancingState = {
    routes: [],
    groups: {},
    planUsageDetails: {
        planName: '',
        computeUsage: null,
        tableUsage: {
            monthlyBudgetTableCount: 0,
            endpointCount: 0,
            tableCount: 0
        }
    },
    tableMetrics: {},
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
        setTableMetrics: (state, action: PayloadAction<any>) => {
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
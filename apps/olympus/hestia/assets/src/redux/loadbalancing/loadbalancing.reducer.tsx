import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Groups, LoadBalancingState} from "./loadbalancing.types";

const initialState: LoadBalancingState = {
    routes: [],
    groups: {},
    planUsageDetails: {},
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
        setUserPlanDetails: (state, action: PayloadAction<any>) => {
            state.planUsageDetails = action.payload;
        },
    }
});

export const {
    setEndpoints,
    setGroupEndpoints,
    setUserPlanDetails} = loadBalancingSlice.actions;
export default loadBalancingSlice.reducer;
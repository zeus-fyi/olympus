import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Groups, LoadBalancingState} from "./loadbalancing.types";

const initialState: LoadBalancingState = {
    routes: [],
    groups: {},
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
    }
});

export const { setEndpoints, setGroupEndpoints} = loadBalancingSlice.actions;
export default loadBalancingSlice.reducer;
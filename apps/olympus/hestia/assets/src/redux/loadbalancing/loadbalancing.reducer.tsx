import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {LoadBalancingRoutesState} from "./loadbalancing.types";

const initialState: LoadBalancingRoutesState = {
    endpoints: [],
}

const loadBalancingSlice = createSlice({
    name: 'loadBalancing',
    initialState,
    reducers: {
        setEndpoints: (state, action: PayloadAction<[any]>) => {
            state.endpoints = action.payload;
        },

    }
});

export const { setEndpoints} = loadBalancingSlice.actions;
export default loadBalancingSlice.reducer;
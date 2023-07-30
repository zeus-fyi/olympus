import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {LoadBalancingState} from "./loadbalancing.types";

const initialState: LoadBalancingState = {
    routes: [],
}

const loadBalancingSlice = createSlice({
    name: 'loadBalancing',
    initialState,
    reducers: {
        setEndpoints: (state, action: PayloadAction<any[]>) => {
            state.routes = action.payload;
        },

    }
});

export const { setEndpoints} = loadBalancingSlice.actions;
export default loadBalancingSlice.reducer;
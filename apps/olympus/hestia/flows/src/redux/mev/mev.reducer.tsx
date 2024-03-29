import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {MevState} from "./mev.actions";

const initialState: MevState = {
    bundles: [],
    callBundles:[],
}

const mevSlice = createSlice({
    name: 'mev',
    initialState,
    reducers: {
        // Action to set the bundles state with an array of any type
        setBundlesState: (state, action: PayloadAction<any[]>) => {
            state.bundles = action.payload;
        },
        setCallBundlesState: (state, action: PayloadAction<[]>) => {
            state.callBundles = action.payload;
        },
    },
});

export const { setBundlesState,setCallBundlesState } = mevSlice.actions;

export default mevSlice.reducer;
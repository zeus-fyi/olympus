import {AppsState, TopologySystemComponentsSlice} from "./apps.types";
import {createSlice, PayloadAction} from "@reduxjs/toolkit";

const initialState: AppsState = {
    privateOrgApps: [],
}

const appsSlice = createSlice({
    name: 'apps',
    initialState,
    reducers: {
        setPrivateOrgApps: (state, action: PayloadAction<TopologySystemComponentsSlice>) => {
            state.privateOrgApps = action.payload;
        },
    }
});

export const { setPrivateOrgApps } = appsSlice.actions;
export default appsSlice.reducer;
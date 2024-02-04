import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {ClustersConfigsState, ExtClusterConfig} from "./clusters.configs.types";

const initialState: ClustersConfigsState = {
   clusterConfigs: [],
}

const clustersConfigsSlice = createSlice({
    name: 'clustersConfigs',
    initialState,
    reducers: {
        // Action to set the bundles state with an array of any type
        setExtClustersConfigs: (state, action: PayloadAction<ExtClusterConfig[]>) => {
            state.clusterConfigs = action.payload;
        },
    },
});

export const { setExtClustersConfigs } = clustersConfigsSlice.actions;

export default clustersConfigsSlice.reducer;
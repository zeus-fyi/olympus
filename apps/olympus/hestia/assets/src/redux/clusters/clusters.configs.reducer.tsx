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
        updateExtClusterConfig: (state, action: PayloadAction<{ index: number; changes: Partial<ExtClusterConfig> }>) => {
            const { index, changes } = action.payload;
            if (state.clusterConfigs[index]) {
                state.clusterConfigs[index] = { ...state.clusterConfigs[index], ...changes };
            }
        },
    },
});

export const { setExtClustersConfigs, updateExtClusterConfig } = clustersConfigsSlice.actions;

export default clustersConfigsSlice.reducer;
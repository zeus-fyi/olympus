import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {ClusterConfig, ClustersConfigsState} from "./clusters.configs.types";

const initialState: ClustersConfigsState = {
   clusterConfigs: [],
}

const clustersConfigsSlice = createSlice({
    name: 'clustersConfigs',
    initialState,
    reducers: {
        // Action to set the bundles state with an array of any type
        setClustersConfigs: (state, action: PayloadAction<ClusterConfig[]>) => {
            state.clusterConfigs = action.payload;
        },
        updateClusterConfigs: (state, action: PayloadAction<{ index: number; changes: Partial<ClusterConfig> }>) => {
            const { index, changes } = action.payload;
            if (state.clusterConfigs[index]) {
                state.clusterConfigs[index] = { ...state.clusterConfigs[index], ...changes };
            }
        },
    },
});

export const { setClustersConfigs, updateClusterConfigs } = clustersConfigsSlice.actions;

export default clustersConfigsSlice.reducer;
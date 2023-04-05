import {AppsState, TopologySystemComponentsSlice} from "./apps.types";
import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ComponentBases, Ingress, IngressPaths} from "../clusters/clusters.types";

const initialState: AppsState = {
    privateOrgApps: [],
    selectedClusterApp: {
        clusterName: '',
        componentBases: {} as ComponentBases,
        ingressSettings: {authServerURL: 'aegis.zeus.fyi', host: 'host.zeus.fyi'} as Ingress,
        ingressPaths: {} as IngressPaths,
    } as Cluster,
    selectedComponentBaseName: '',
    selectedSkeletonBaseName: '',
}

const appsSlice = createSlice({
    name: 'apps',
    initialState,
    reducers: {
        setPrivateOrgApps: (state, action: PayloadAction<TopologySystemComponentsSlice>) => {
            state.privateOrgApps = action.payload;
        },
        setSelectedClusterApp: (state, action: PayloadAction<Cluster>) => {
            state.selectedClusterApp = action.payload;
        },
        setSelectedComponentBaseName: (state, action: PayloadAction<string>) => {
            state.selectedComponentBaseName = action.payload;
        },
        setSelectedSkeletonBaseName: (state, action: PayloadAction<string>) => {
            state.selectedSkeletonBaseName = action.payload;
        },
    }
});

export const { setPrivateOrgApps, setSelectedClusterApp, setSelectedSkeletonBaseName, setSelectedComponentBaseName } = appsSlice.actions;
export default appsSlice.reducer;
import {AppsState, TopologySystemComponentsSlice} from "./apps.types";
import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ClusterPreview, ComponentBases, Ingress, IngressPaths} from "../clusters/clusters.types";

const initialState: AppsState = {
    privateOrgApps: [],
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
        ingressSettings: {authServerURL: 'aegis.zeus.fyi', host: 'host.zeus.fyi'} as Ingress,
        ingressPaths: {} as IngressPaths,
    } as Cluster,
    clusterPreview: {} as ClusterPreview,
    selectedComponentBaseName: '',
    selectedSkeletonBaseName: '',
}

const appsSlice = createSlice({
    name: 'apps',
    initialState,
    reducers: {
        setClusterPreview: (state, action: PayloadAction<ClusterPreview>) => {
            state.clusterPreview = action.payload;
        },
        setPrivateOrgApps: (state, action: PayloadAction<TopologySystemComponentsSlice>) => {
            state.privateOrgApps = action.payload;
        },
        setCluster: (state, action: PayloadAction<Cluster>) => {
            state.cluster = action.payload;
        },
        setSelectedComponentBaseName: (state, action: PayloadAction<string>) => {
            state.selectedComponentBaseName = action.payload;
        },
        setSelectedSkeletonBaseName: (state, action: PayloadAction<string>) => {
            state.selectedSkeletonBaseName = action.payload;
        },
    }
});

export const { setPrivateOrgApps,setClusterPreview, setCluster, setSelectedSkeletonBaseName, setSelectedComponentBaseName } = appsSlice.actions;
export default appsSlice.reducer;
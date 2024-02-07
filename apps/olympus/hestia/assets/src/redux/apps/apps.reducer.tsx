import {AppsState, Nodes, TopologySystemComponentsSlice} from "./apps.types";
import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ClusterPreview, ComponentBases, Ingress, IngressPaths} from "../clusters/clusters.types";
import {CloudProviderRegionsResourcesMap} from "../resources/resources.types";

const initialState: AppsState = {
    privateOrgApps: [],
    publicMatrixFamilyApps: [],
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
        ingressSettings: {authServerURL: 'aegis.zeus.fyi', host: 'host.zeus.fyi'} as Ingress,
        ingressPaths: {} as IngressPaths,
    } as Cluster,
    clusterPreview: {} as ClusterPreview,
    selectedComponentBaseName: '',
    selectedSkeletonBaseName: '',
    nodes: [{
        resourceID: 0,
        description: '',
        slug: '',
        disk: 0,
        priceHourly: 0,
        cloudProvider: '',
        vcpus: 0,
        priceMonthly: 0,
        region: '',
        memory: 0,
        gpus: 0,
        gpuType: 'none',
    }],
    cloudRegionResourceMap: {} as CloudProviderRegionsResourcesMap,
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
        setPublicMatrixFamilyApps: (state, action: PayloadAction<TopologySystemComponentsSlice>) => {
            state.publicMatrixFamilyApps = action.payload;
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
        setNodes: (state, action: PayloadAction<Nodes[]>) => {
            state.nodes = action.payload;
        },
        setCloudRegionResourceMap: (state, action: PayloadAction<CloudProviderRegionsResourcesMap>) => {
            state.cloudRegionResourceMap = action.payload;
        },
    }
});

export const { setPublicMatrixFamilyApps,
    setNodes,
    setPrivateOrgApps,
    setCloudRegionResourceMap,
    setClusterPreview,
    setCluster,
    setSelectedSkeletonBaseName,
    setSelectedComponentBaseName } = appsSlice.actions;
export default appsSlice.reducer;
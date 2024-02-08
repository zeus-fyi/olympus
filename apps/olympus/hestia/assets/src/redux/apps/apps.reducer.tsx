import {AppsState, Nodes, TopologySystemComponentsSlice} from "./apps.types";
import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ClusterPreview, ComponentBases, Ingress, IngressPaths} from "../clusters/clusters.types";
import {CloudProviderRegionsResourcesMap, Disks} from "../resources/resources.types";

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
    selectedCloudProvider: '',
    selectedRegion: '',
    selectedDisk: {
        extCfgStrID: '',
        resourceStrID: '',
        diskUnits: '',
        priceMonthly: 0,
        description: '',
        type: '',
        subType: '',
        diskSize: 0,
        priceHourly: 0,
        region: '',
        cloudProvider: '',
    },
    selectedNode: {
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
    },
    deployServersCount: 0,
    cloudRegionResourceMap: {} as CloudProviderRegionsResourcesMap,
}

const appsSlice = createSlice({
    name: 'apps',
    initialState,
    reducers: {
        setDeployServersCount: (state, action: PayloadAction<number>) => {
            state.deployServersCount = action.payload;
        },
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
        setCloudProvider: (state, action: PayloadAction<string>) => {
            state.selectedCloudProvider = action.payload;
        },
        setRegion: (state, action: PayloadAction<string>) => {
            state.selectedRegion = action.payload;
        },
        setSelectedDisk: (state, action: PayloadAction<Disks>) => {
            state.selectedDisk = action.payload;
        },
        setSelectedNode: (state, action: PayloadAction<Nodes>) => {
            state.selectedNode = action.payload;
        },
        setCloudRegionResourceMap: (state, action: PayloadAction<CloudProviderRegionsResourcesMap>) => {
            state.cloudRegionResourceMap = action.payload;

            const cloudProviderKeys = Object.keys(action.payload).filter(key => key.trim() !== '');
            if (cloudProviderKeys.length > 0) {
                const firstValidCloudProvider = cloudProviderKeys[0];
                state.selectedCloudProvider = firstValidCloudProvider;

                // Now find the first valid region for the selected cloud provider
                const regionKeys = Object.keys(action.payload[firstValidCloudProvider]).filter(key => key.trim() !== '');
                if (regionKeys.length > 0) {
                    state.selectedRegion = regionKeys[0];


                    // Now find the first valid node for the selected region
                    const nodes = action.payload[firstValidCloudProvider][regionKeys[0]].nodes;
                    if (nodes && nodes.length > 0) {
                        state.selectedNode = nodes[0];
                    } else {
                        // If no valid nodes found, set to an empty object or a default value
                        state.selectedNode = initialState.selectedNode;
                    }

                    const disks = action.payload[firstValidCloudProvider][regionKeys[0]].disks;
                    if (disks.length > 0) {
                        state.selectedDisk = disks[0];
                    } else {
                        // If no valid disks found, set to an empty object or a default value
                        state.selectedDisk = {
                            extCfgStrID: '',
                            resourceStrID: '',
                            diskUnits: '',
                            priceMonthly: 0,
                            description: '',
                            type: '',
                            subType: '',
                            diskSize: 0,
                            priceHourly: 0,
                            region: '',
                            cloudProvider: '',
                        };
                    }

                } else {
                    // If no valid region found, set to an empty string or a default value
                    state.selectedRegion = '';
                }
            } else {
                // If no valid cloud provider found, set both to empty strings or default values
                state.selectedCloudProvider = '';
                state.selectedRegion = '';
            }
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
    setSelectedComponentBaseName,
    setCloudProvider,
    setRegion,
    setSelectedDisk,
    setSelectedNode,
    setDeployServersCount
} = appsSlice.actions;
export default appsSlice.reducer;
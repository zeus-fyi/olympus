import {Cluster} from "../clusters/clusters.types";

export interface TopologySystemComponents {
    topologySystemComponentID: number;
    topologyClassTypeID: number;
    topologySystemComponentName: string;
}

export interface AppsState {
    privateOrgApps: TopologySystemComponentsSlice;
    selectedClusterApp: Cluster;
}

export type TopologySystemComponentsSlice = TopologySystemComponents[];
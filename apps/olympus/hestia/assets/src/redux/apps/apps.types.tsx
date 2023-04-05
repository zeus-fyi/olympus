import {Cluster} from "../clusters/clusters.types";

export interface TopologySystemComponents {
    topologySystemComponentID: number;
    topologyClassTypeID: number;
    topologySystemComponentName: string;
}

export interface AppsState {
    privateOrgApps: TopologySystemComponentsSlice;
    selectedClusterApp: Cluster;
    selectedComponentBaseName: string;
    selectedSkeletonBaseName: string;
}

export type TopologySystemComponentsSlice = TopologySystemComponents[];
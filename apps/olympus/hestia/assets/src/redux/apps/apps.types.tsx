import {Cluster, ClusterPreview} from "../clusters/clusters.types";

export interface TopologySystemComponents {
    topologySystemComponentID: number;
    topologyClassTypeID: number;
    topologySystemComponentName: string;
}

export interface AppsState {
    privateOrgApps: TopologySystemComponentsSlice;
    cluster: Cluster;
    clusterPreview: ClusterPreview;
    selectedComponentBaseName: string;
    selectedSkeletonBaseName: string;
}

export type TopologySystemComponentsSlice = TopologySystemComponents[];
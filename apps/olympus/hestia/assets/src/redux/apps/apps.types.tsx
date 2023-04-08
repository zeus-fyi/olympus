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
    nodes: Nodes[];
}

export interface Nodes {
    nodeID: number;
    description: string;
    slug: string;
    disk: number;
    priceHourly: number;
    cloudProvider: string;
    vcpus: number;
    priceMonthly: number;
    region: string;
    memory: number;
}

export type TopologySystemComponentsSlice = TopologySystemComponents[];
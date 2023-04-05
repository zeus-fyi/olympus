export interface TopologySystemComponents {
    topologySystemComponentID: number;
    topologyClassTypeID: number;
    topologySystemComponentName: string;
}

export interface AppsState {
    privateOrgApps: TopologySystemComponentsSlice;
}

export type TopologySystemComponentsSlice = TopologySystemComponents[];
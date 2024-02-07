export interface ClustersConfigsState {
    clusterConfigs: ClusterConfig[];
    authedClusterConfigs: ClusterConfig[];
}

export interface ClusterConfig{
    extConfigStrID: string,
    cloudProvider: string,
    region: string,
    context: string,
    contextAlias: string
    env: string,
    isActive: boolean
}


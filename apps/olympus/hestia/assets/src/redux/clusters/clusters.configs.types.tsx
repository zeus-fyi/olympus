

export interface ClustersConfigsState {
    clusterConfigs: ExtClusterConfig[];
}


export interface ExtClusterConfig{
    extConfigStrID: string,
    cloudProvider: string,
    region: string,
    context: string,
    contextAlias: string
    env: string,
    isActive: boolean
}

export interface ClustersConfigsState {
    clusterConfigs: ClusterConfig[];
    authedClusterConfigs: ClusterConfig[];
}

export interface ClusterConfig{
    extConfigStrID: string,
    cloudCtxNs: CloudCtxNs,
    contextAlias: string
    isActive: boolean
}

interface CloudCtxNs {
    clusterCfgStrID: string; // The `?` indicates that this field is optional, aligning with the `omitempty` in the Go struct tag.
    cloudProvider: string;
    region: string;
    context: string;
    namespace: string;
    alias: string; // Optional field, similar to `clusterCfgStrID`.
    env: string; // Optional field.
}

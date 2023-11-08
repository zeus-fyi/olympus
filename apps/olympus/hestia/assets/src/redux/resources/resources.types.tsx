import {NodeStatus, V1Taint} from "@kubernetes/client-node";

export interface ResourcesState {
    resources: any[];
    searchResources: any[];
    appNodes: NodeAudit[];
}

export interface NodeAudit {
    kubernetesVersion: string;
    nodeID: string;
    nodePoolID: string;
    region: string;
    slug: string;
    taints: V1Taint[];
    status: NodeStatus;
}
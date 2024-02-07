import {NodeStatus, V1Taint} from "@kubernetes/client-node";

export interface ResourcesState {
    resources: any[];
    searchResources: NodesSlice;
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

type CloudProviderRegionsMap = { [key: string]: string[] };

export interface ResourceAggregate {
    monthlyPrice?: number;
    hourlyPrice?: number;
    memRequests: string;
    cpuRequests: string;
}

export interface ResourceMinMax {
    max: ResourceAggregate;
    min: ResourceAggregate;
}

export interface NodeSearchParams {
    cloudProviderRegions: CloudProviderRegionsMap;
    diskType?: string;
    resourceMinMax?: ResourceMinMax;
}

export interface NodeSearchRequest {
    nodeSearchParams: NodeSearchParams;
}

export interface Node {
    extCfgStrID: string;
    memory: number;
    vcpus: number;
    disk: number;
    diskUnits: string;
    diskType: string;
    priceHourly: number;
    region: string;
    cloudProvider: string;
    resourceID: number;
    description: string;
    slug: string;
    memoryUnits: string;
    priceMonthly: number;
    gpus: number;
    gpuType: string;
}

export type NodesSlice = Node[];

export interface RegionResourcesMap {
    [region: string]: Resources;
}

export interface Resources {
    nodes: Node[];
    disks: Disks[];
}

export interface Disks {
    extCfgStrID: string;

    resourceID: number;
    diskUnits: string;
    priceMonthly: number;
    description: string;
    type: string;
    diskSize: number;
    priceHourly: number;
    region: string;
    cloudProvider: string;
}


export interface CloudProviderRegionsResourcesMap {
    [provider: string]: RegionResourcesMap;
}


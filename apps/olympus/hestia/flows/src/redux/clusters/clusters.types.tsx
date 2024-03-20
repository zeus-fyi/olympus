import {V1ConfigMap, V1Deployment, V1Ingress, V1Service, V1StatefulSet} from '@kubernetes/client-node';

export interface Cluster{
    clusterName: string;
    componentBases: ComponentBases;
    ingressSettings: Ingress;
    ingressPaths: IngressPaths;
}

export interface ClusterPreview {
    clusterName: string;
    componentBases: ComponentBasesPreviews
}

export type ComponentBasesPreviews = {
    [componentBaseName: string]: SkeletonBasesPreviews;
};

export type SkeletonBasesPreviews = {
    [skeletonBaseName: string]: SkeletonBasePreview;
};

export type SkeletonBasePreview = {
    service: V1Service | null;
    configMap: V1ConfigMap | null;
    deployment: V1Deployment | null;
    statefulSet: V1StatefulSet | null;
    ingress: V1Ingress | null;
}

export type IngressPaths = {
    [componentBaseName: string]: IngressPath;
};

export type ComponentBases = {
    [componentBaseName: string]: SkeletonBases;
};

export type SkeletonBases = {
    [skeletonBaseName: string]: SkeletonBase;
};

export interface SkeletonBase {
    addStatefulSet: boolean;
    addDeployment: boolean;
    addConfigMap: boolean;
    addService: boolean;
    addIngress: boolean;
    addServiceMonitor: boolean;
    configMap: ConfigMap;
    deployment: Deployment;
    statefulSet: StatefulSet;
    containers: Containers;
    resourceSums?: ResourceSums;
}

export interface ResourceSums {
    memRequests: string;
    memLimits: string;
    cpuRequests: string;
    cpuLimits: string;
    diskRequests: string;
    diskLimits: string;
    replicas: string;
}

export interface ConfigMap {
    [key: string]: string;
}

export interface Ingress {
    authServerURL: string
    host: string
}

export interface IngressPath {
    path: string
    pathType: string
}

export interface Deployment {
    replicaCount: number;
}

export interface StatefulSet {
    replicaCount: number;
    pvcTemplates: PVCTemplate[];
}

export interface PVCTemplate {
    name: string;
    accessMode: string;
    storageSizeRequest: string;
}

// just conditionally add/remove items if deployment or stateful set is selected
export type Containers = {
    [containerName: string]: Container;
};

export interface Container {
    isInitContainer: boolean;
    dockerImage: DockerImage;
}

export interface DockerImage {
    imageName: string;
    cmd: string;
    args: string;
    resourceRequirements: ResourceRequirements;
    ports: Port[];
    volumeMounts: VolumeMount[];
}

export interface ResourceRequirements {
    cpu: string;
    memory: string;
}

export interface Port {
    name: string;
    number: number;
    protocol: string;
    ingressEnabledPort: boolean;
}

export interface VolumeMount {
    name: string;
    mountPath: string;
}

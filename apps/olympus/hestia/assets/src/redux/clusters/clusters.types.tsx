export interface Cluster{
    clusterName: string;
    componentBases: ComponentBases;
}

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
    ingress: Ingress;
    deployment: Deployment;
    statefulSet: StatefulSet;
    containers: Containers;
}

export interface Ingress {
    authServerURL: string
    host: string
    paths: IngressPath[]
}

export interface IngressPath {
    path: string
    pathType: string
    serviceName: string
}
export interface ConfigMap {
    [key: string]: string;
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
}

export interface VolumeMount {
    name: string;
    mountPath: string;
}
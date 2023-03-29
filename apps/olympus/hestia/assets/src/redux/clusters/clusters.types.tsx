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
    containers: Containers;
}

export type Containers = {
    [containerName: string]: Container;
};

export interface Container {
    dockerImage: DockerImage;
}

export interface DockerImage {
    imageName: string;
    cmd: string;
    args: string[];
    ports: Port[];
}

export interface Port {
    name: string;
    number: number;
    protocol: string;
}
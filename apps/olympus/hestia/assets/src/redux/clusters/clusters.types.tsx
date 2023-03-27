
export type Cluster = {
    clusterName: string;
    componentBases: ComponentBases;
};

export type ComponentBases = {
    [key: string]: SkeletonBases;
};

export type SkeletonBases = {
    [key: string]: SkeletonBase;
};

export interface SkeletonBase {
    dockerImages: DockerImages;
}

export type DockerImages = {
    [key: string]: DockerImage;
};

export interface DockerImage {
    imageName: string;
    ports: Port[];
}

export interface Port {
    name: string;
    number: number;
    protocol: string;
}
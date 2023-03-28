import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ComponentBases, DockerImage, Port, SkeletonBase, SkeletonBases} from "./clusters.types";

interface ClusterBuilderState {
    cluster: Cluster;
    selectedComponentBase: SkeletonBases;
    selectedSkeletonBase: SkeletonBase;
}

const initialState: ClusterBuilderState = {
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
    },
    selectedComponentBase: {} as SkeletonBases,
    selectedSkeletonBase: {} as SkeletonBase,
};

const clusterBuilderSlice = createSlice({
    name: 'clusterBuilder',
    initialState,
    reducers: {
        setClusterName: (state, action: PayloadAction<string>) => {
            state.cluster.clusterName = action.payload;
        },
        addComponentBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBases: SkeletonBases }>) => {
            const { componentBaseName, skeletonBases } = action.payload;
            state.cluster.componentBases[componentBaseName] = skeletonBases;
        },
        removeComponentBase: (state, action: PayloadAction<string>) => {
            const key = action.payload;
            if (state.cluster.componentBases[key]) {
                delete state.cluster.componentBases[key];
            } else {
                console.error(`Component base not found: ${key}`);
            }
        },
        addSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; skeletonBase: SkeletonBase }>) => {
            const { componentBaseName, skeletonBaseName, skeletonBase } = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName] = skeletonBase;
        },
        removeSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string }>) => {
            const { componentBaseName, skeletonBaseName } = action.payload;
            if (state.cluster.componentBases[componentBaseName][skeletonBaseName]) {
                delete state.cluster.componentBases[componentBaseName][skeletonBaseName];
            } else {
                console.error(`Skeleton base not found: ${skeletonBaseName}`);
            }
        },
        addDockerImage: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; key: string; dockerImage: DockerImage }>) => {
            const { componentBaseKey, skeletonBaseKey, key, dockerImage } = action.payload;
            if (!state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]) {
                console.error(`SkeletonBase not found: ${skeletonBaseKey}`);
                return;
            }
            state.cluster.componentBases[componentBaseKey][skeletonBaseKey].dockerImages[key] = dockerImage;
        },
        setDockerImagePort: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; dockerImageKey: string; portIndex: number; port: Port }>) => {
            const { componentBaseKey, skeletonBaseKey, dockerImageKey, portIndex, port } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.dockerImages[dockerImageKey];
            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            if (portIndex < 0 || portIndex >= dockerImage.ports.length) {
                console.error(`Invalid port index: ${portIndex}`);
                return;
            }
            dockerImage.ports[portIndex] = port;
        },
    },
});

export const { setClusterName, addComponentBase, removeComponentBase, addSkeletonBase, removeSkeletonBase, addDockerImage, setDockerImagePort} = clusterBuilderSlice.actions;

export default clusterBuilderSlice.reducer;

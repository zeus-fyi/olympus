import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ComponentBases, DockerImage, Port, SkeletonBase, SkeletonBases} from "./clusters.types";

interface ClusterBuilderState {
    cluster: Cluster;
}

const initialState: ClusterBuilderState = {
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
    },
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
        addSkeletonBase: (state, action: PayloadAction<{ componentBaseKey: string; key: string; skeletonBase: SkeletonBase }>) => {
            const { componentBaseKey, key, skeletonBase } = action.payload;
            if (!state.cluster.componentBases[componentBaseKey]) {
                state.cluster.componentBases[componentBaseKey] = {};
            }
            state.cluster.componentBases[componentBaseKey][key] = skeletonBase;
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

export const { setClusterName, addComponentBase, removeComponentBase, addSkeletonBase, addDockerImage, setDockerImagePort} = clusterBuilderSlice.actions;

export default clusterBuilderSlice.reducer;

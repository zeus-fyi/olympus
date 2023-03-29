import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Cluster, ComponentBases, Container, Port, SkeletonBase, SkeletonBases} from "./clusters.types";

interface ClusterBuilderState {
    cluster: Cluster;
    selectedComponentBaseName: string;
    selectedSkeletonBaseName: string;
    selectedDockerImageName: string;
}

const initialState: ClusterBuilderState = {
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
    },
    selectedComponentBaseName: '',
    selectedSkeletonBaseName: '',
    selectedDockerImageName: '',
};

const clusterBuilderSlice = createSlice({
    name: 'clusterBuilder',
    initialState,
    reducers: {
        setClusterName: (state, action: PayloadAction<string>) => {
            state.cluster.clusterName = action.payload;
        },
        setSelectedDockerImageName: (state, action: PayloadAction<string>) => {
            state.selectedDockerImageName = action.payload;
        },
        setSelectedComponentBaseName: (state, action: PayloadAction<string>) => {
            state.selectedComponentBaseName = action.payload;
        },
        setSelectedSkeletonBaseName: (state, action: PayloadAction<string>) => {
            state.selectedSkeletonBaseName = action.payload;
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
        addContainer: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; container: Container }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, container } = action.payload;
            if (!state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]) {
                console.error(`SkeletonBase not found: ${skeletonBaseKey}`);
                return;
            }
            state.cluster.componentBases[componentBaseKey][skeletonBaseKey].containers[containerName] = container;
        },
        setDockerImagePort: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; dockerImageKey: string; portIndex: number; port: Port }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, portIndex, port } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
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

export const { setClusterName, addComponentBase, removeComponentBase, addSkeletonBase,
    setSelectedDockerImageName, removeSkeletonBase, setSelectedComponentBaseName,setSelectedSkeletonBaseName,
    addContainer, setDockerImagePort} = clusterBuilderSlice.actions;

export default clusterBuilderSlice.reducer;

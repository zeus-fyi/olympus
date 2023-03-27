import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {DockerImages, Port} from "./clusters.types";

interface ClusterBuilderState {
    clusterName: string;
    dockerImages: DockerImages;
}

const initialState: ClusterBuilderState = {
    clusterName: '',
    dockerImages: {},
};

const clusterBuilderSlice = createSlice({
    name: 'clusterBuilder',
    initialState,
    reducers: {
        setClusterName: (state, action: PayloadAction<string>) => {
            state.clusterName = action.payload;
        },
        addDockerImage: (state, action: PayloadAction<any>) => {
            state.dockerImages.dockerImageName = action.payload;
        },
        addDockerImagePort: (state, action: PayloadAction<{dockerImageName: string, port: Port}>) => {
            const { dockerImageName, port } = action.payload;

            if (state.dockerImages[dockerImageName]) {
                state.dockerImages[dockerImageName].ports.push(port);
            } else {
                console.error(`Docker image not found: ${dockerImageName}`);
            }
        },
    },
});

export const { setClusterName, addDockerImagePort} = clusterBuilderSlice.actions;

export default clusterBuilderSlice.reducer;

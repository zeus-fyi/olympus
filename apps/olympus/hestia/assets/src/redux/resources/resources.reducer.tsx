import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {NodeAudit, ResourcesState} from "./resources.types";

const initialState: ResourcesState = {
    resources: [],
    searchResources: [],
    appNodes: [],
}

const resourcesSlice = createSlice({
    name: 'resources',
    initialState,
    reducers: {
        setResources: (state, action: PayloadAction<[any]>) => {
            state.resources = action.payload;
        },
        setSearchResources: (state, action: PayloadAction<[any]>) => {
            state.searchResources = action.payload;
        },
        setAppNodes: (state, action: PayloadAction<NodeAudit[]>) => {
            state.appNodes = action.payload;
        }
    }
});

export const { setResources, setAppNodes, setSearchResources } = resourcesSlice.actions;
export default resourcesSlice.reducer;
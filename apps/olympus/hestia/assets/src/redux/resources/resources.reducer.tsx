import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {NodeAudit, ResourcesState} from "./resources.types";

const initialState: ResourcesState = {
    resources: [],
    appNodes: [],
}

const resourcesSlice = createSlice({
    name: 'resources',
    initialState,
    reducers: {
        setResources: (state, action: PayloadAction<[any]>) => {
            state.resources = action.payload;
        },
        setAppNodes: (state, action: PayloadAction<NodeAudit[]>) => {
            state.appNodes = action.payload;
        }
    }
});

export const { setResources, setAppNodes } = resourcesSlice.actions;
export default resourcesSlice.reducer;
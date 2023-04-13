import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {ResourcesState} from "./resources.types";

const initialState: ResourcesState = {
    resources: [],
}

const resourcesSlice = createSlice({
    name: 'resources',
    initialState,
    reducers: {
        setResources: (state, action: PayloadAction<[any]>) => {
            state.resources = action.payload;
        },
    }
});

export const { setResources } = resourcesSlice.actions;
export default resourcesSlice.reducer;
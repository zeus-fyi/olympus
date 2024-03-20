import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {initialState} from "./flows.actions";

const flowsSlice = createSlice({
    name: 'flows',
    initialState,
    reducers: {
        setUploadContent: (state, action: PayloadAction<any>) => {
            state.uploadContent = action.payload;
        },
    }
});

export const {
    setUploadContent
} = flowsSlice.actions;
export default flowsSlice.reducer
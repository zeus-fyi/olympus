import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {initialState} from "./flows.actions";

const flowsSlice = createSlice({
    name: 'flows',
    initialState,
    reducers: {
        setUploadContent: (state, action: PayloadAction<any>) => {
            state.uploadContent = action.payload;
        },
        setCsvHeaders: (state, action: PayloadAction<string[]>) => {
            state.csvHeaders = action.payload;
        }
    }
});

export const {
    setCsvHeaders,
    setUploadContent
} = flowsSlice.actions;
export default flowsSlice.reducer
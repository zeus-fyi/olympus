import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {initialState} from "./flows.actions";

const flowsSlice = createSlice({
    name: 'flows',
    initialState,
    reducers: {
        setUploadTasksContent: (state, action: PayloadAction<any>) => {
            state.uploadContentTasks = action.payload;
        },
        setUploadContacts: (state, action: PayloadAction<any>) => {
            state.uploadContentContacts = action.payload;
        },
        setCsvHeaders: (state, action: PayloadAction<string[]>) => {
            state.csvHeaders = action.payload;
        },
        setPromptHeaders: (state, action: PayloadAction<string[]>) => {
            state.promptHeaders = action.payload;
        },
        setResults: (state, action: PayloadAction<any[]>) => {
            state.results = action.payload;
        }
    }
});
export const {
    setCsvHeaders,
    setUploadContacts,
    setUploadTasksContent,
    setPromptHeaders,
    setResults
} = flowsSlice.actions;
export default flowsSlice.reducer
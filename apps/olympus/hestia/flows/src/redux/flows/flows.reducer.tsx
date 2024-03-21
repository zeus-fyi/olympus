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
        setResults: (state, action: PayloadAction<[]>) => {
            state.results = action.payload;
        },
        setStages: (state, action: PayloadAction<[]>) => {
            state.stages = action.payload;
        }
    }
});

export const {
    setCsvHeaders,
    setUploadContacts,
    setUploadTasksContent,
    setPromptHeaders,
    setResults,
    setStages
} = flowsSlice.actions;
export default flowsSlice.reducer
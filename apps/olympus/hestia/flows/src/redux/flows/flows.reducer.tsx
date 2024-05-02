import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {initialState} from "./flows.actions";

const flowsSlice = createSlice({
    name: 'flows',
    initialState,
    reducers: {
        setPromptsCsvContent: (state, action: PayloadAction<[]>) => {
            state.promptsCsvContent = action.payload;
        },
        setUploadContacts: (state, action: PayloadAction<[]>) => {
            state.uploadContentContacts = action.payload;
        },
        setCsvHeaders: (state, action: PayloadAction<string[]>) => {
            state.csvHeaders = action.payload;
        },
        setPreviewCount: (state, action: PayloadAction<number>) => {
            state.previewCount = action.payload;
        },
        setPromptHeaders: (state, action: PayloadAction<string[]>) => {
            state.promptHeaders = action.payload;
        },
        setResults: (state, action: PayloadAction<[]>) => {
            state.results = action.payload;
        },
        setStages: (state, action: PayloadAction<{ [key: string]: boolean }>) => {
            state.stages = {
                ...state.stages,
                ...action.payload
            };
        },
        setCommandPrompt: (state, action: PayloadAction<{ [key: string]: string }>) => {
            state.commandPrompts = {
                ...state.commandPrompts,
                ...action.payload
            };
        }
    }
});

export const {
    setCsvHeaders,
    setUploadContacts,
    setPromptsCsvContent,
    setPromptHeaders,
    setResults,
    setStages,
    setCommandPrompt,
    setPreviewCount
} = flowsSlice.actions;
export default flowsSlice.reducer
import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {initialState, UpdateTaskRelationshipPayload} from "./flows.actions";

const flowsSlice = createSlice({
    name: 'flows',
    initialState,
    reducers: {
        setFlowList: (state, action: PayloadAction<string[]>) => {
            state.flowList = action.payload;
        },
        setPromptsCsvContent: (state, action: PayloadAction<[]>) => {
            state.promptsCsvContent = action.payload;
        },
        setUploadContacts: (state, action: PayloadAction<[]>) => {
            state.uploadContentContacts = action.payload;
        },
        setCsvHeaders: (state, action: PayloadAction<string[]>) => {
            state.csvHeaders = action.payload;
        },
        setContactsCsvFilename: (state, action: PayloadAction<string>) => {
            state.contactsCsvFilename = action.payload;
        },
        setPreviewCount: (state, action: PayloadAction<number>) => {
            state.previewCount = action.payload;
        },
        setPromptHeaders: (state, action: PayloadAction<string[]>) => {
            state.promptHeaders = action.payload;
            state.stagePromptMap = {}
            action.payload.forEach(header => {
                state.stagePromptMap[header] = "Default";
            });
        },
        setPromptFlowMap: (state, action: PayloadAction<UpdateTaskRelationshipPayload>) => {
            const { key, subKey } = action.payload;
            state.stagePromptMap[key] = subKey
        },
        setContactsFlowMap: (state, action: PayloadAction<UpdateTaskRelationshipPayload>) => {
            const { key, subKey } = action.payload;
            state.stageContactsMap[key] = subKey;
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
    setPreviewCount,
    setContactsCsvFilename,
    setPromptFlowMap,
    setContactsFlowMap,
    setFlowList
} = flowsSlice.actions;
export default flowsSlice.reducer
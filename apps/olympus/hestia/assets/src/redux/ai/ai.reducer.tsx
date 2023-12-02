import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {AiState, TaskModelInstructions, UpdateTaskMapPayload} from "./ai.types";

const initialState: AiState = {
    searchContentText: '',
    usernames: '',
    groupFilter: '',
    analysisWorkflowInstructions: '',
    aggregationWorkflowInstructions: '',
    searchResults: '',
    platformFilter: '',
    workflows: [],
    tasks: [],
    addAnalysisView: false,
    addAggregationView: false,
    addedAnalysisTasks: [],
    addedAggregateTasks: [],
    workflowBuilderTaskMap: {}
}

const aiSlice = createSlice({
    name: 'ai',
    initialState,
    reducers: {
        setAddAnalysisView: (state, action: PayloadAction<boolean>) => {
            state.addAnalysisView = action.payload;
        },
        setAddAggregationView: (state, action: PayloadAction<boolean>) => {
            state.addAggregationView = action.payload;
        },
        setUsernames: (state, action: PayloadAction<string>) => {
            state.usernames = action.payload;
        },
        setGroupFilter: (state, action: PayloadAction<string>) => {
            state.groupFilter = action.payload;
        },
        setSearchContent: (state, action: PayloadAction<string>) => {
            state.searchContentText = action.payload;
        },
        setAnalysisWorkflowInstructions: (state, action: PayloadAction<string>) => {
            state.analysisWorkflowInstructions = action.payload;
        },
        setAggregationWorkflowInstructions: (state, action: PayloadAction<string>) => {
            state.aggregationWorkflowInstructions = action.payload;
        },
        setSearchResults: (state, action: PayloadAction<string>) => {
            state.searchResults = action.payload;
        },
        setPlatformFilter: (state, action: PayloadAction<string>) => {
            state.platformFilter = action.payload;
        },
        setWorkflows: (state, action: PayloadAction<[]>) => {
            state.workflows = action.payload;
        },
        setTasks: (state, action: PayloadAction<[]>) => {
            state.tasks = action.payload;
        },
        setAddAnalysisTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAnalysisTasks = action.payload;
        },
        setAddAggregateTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAggregateTasks = action.payload;
        },
        setWorkflowBuilderTaskMap: (state, action: PayloadAction<UpdateTaskMapPayload>) => {
            const { key, subKey, value } = action.payload;

            if (value) {
                if (!state.workflowBuilderTaskMap[key]) {
                    state.workflowBuilderTaskMap[key] = {};
                }
                state.workflowBuilderTaskMap[key][subKey] = true;
            } else {
                if (state.workflowBuilderTaskMap[key]) {
                    delete state.workflowBuilderTaskMap[key][subKey];

                    // Check if the main key has no inner keys left
                    if (Object.keys(state.workflowBuilderTaskMap[key]).length === 0) {
                        // If so, delete the main key from the map
                        delete state.workflowBuilderTaskMap[key];
                    }
                }
            }
        },
    }
});

export const { setSearchContent,
    setGroupFilter,
    setUsernames,
    setAnalysisWorkflowInstructions,
    setAggregationWorkflowInstructions,
    setSearchResults,
    setWorkflows,
    setPlatformFilter,
    setTasks,
    setAddAnalysisView,
    setAddAggregationView,
    setAddAnalysisTasks,
    setAddAggregateTasks,
    setWorkflowBuilderTaskMap,

} = aiSlice.actions;
export default aiSlice.reducer;
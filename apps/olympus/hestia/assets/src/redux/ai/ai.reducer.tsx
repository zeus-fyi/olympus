import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {AiState, Retrieval, TaskModelInstructions, UpdateTaskCycleCountPayload, UpdateTaskMapPayload} from "./ai.types";

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
    addRetrievalView: false,
    addedAnalysisTasks: [],
    addedAggregateTasks: [],
    addedRetrievals: [],
    workflowBuilderTaskMap: {},
    taskMap: {},
    retrieval: {
        retrievalName: '',
        retrievalGroup: '',
        retrievalKeywords: '',
        retrievalPlatform: '',
        retrievalUsernames: '',
        retrievalPrompt: '',
        retrievalPlatformGroups: '',
    },
    retrievals: [],
    workflowAnalysisRetrievalsMap: {}
}

const aiSlice = createSlice({
    name: 'ai',
    initialState,
    reducers: {
        setRetrievalName: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalName = action.payload;
        },
        setRetrievalPlatformGroups: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalPlatformGroups = action.payload;
        },
        setRetrievalGroup: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalGroup = action.payload;
        },
        setRetrievalKeywords: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalKeywords = action.payload;
        },
        setRetrievalUsernames: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalUsernames = action.payload;
        },
        setRetrievalPlatform: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalPlatform = action.payload;
        },
        setRetrievalPrompt: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalPrompt = action.payload;
        },
        setAddAnalysisView: (state, action: PayloadAction<boolean>) => {
            state.addAnalysisView = action.payload;
        },
        setAddAggregationView: (state, action: PayloadAction<boolean>) => {
            state.addAggregationView = action.payload;
        },
        setAddRetrievalView: (state, action: PayloadAction<boolean>) => {
            state.addRetrievalView = action.payload;
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
        setAiTasks: (state, action: PayloadAction<[]>) => {
            state.tasks = action.payload;
        },
        setRetrievals: (state, action: PayloadAction<[]>) => {
            state.retrievals = action.payload;
        },
        setAddAnalysisTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAnalysisTasks = action.payload;
            for (let i = 0; i < state.addedAnalysisTasks.length; i++) {
                const task  = state.addedAnalysisTasks[i]
                if (task && task.taskID) {
                    state.taskMap[task.taskID] = task;
                }
            }
        },
        setAddAggregateTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAggregateTasks = action.payload;
            for (let i = 0; i < state.addedAggregateTasks.length; i++) {
                const task  = state.addedAggregateTasks[i]
                if (task && task.taskID) {
                    state.taskMap[task.taskID] = task;
                }
            }
        },
        setAddRetrievalTasks: (state, action: PayloadAction<Retrieval[]>) => {
            state.addedRetrievals = action.payload;
        },
        setAnalysisRetrievalsMap: (state, action: PayloadAction<UpdateTaskMapPayload>) => {
            const { key, subKey, value } = action.payload;
            if (value) {
                if (!state.workflowAnalysisRetrievalsMap[key]) {
                    state.workflowAnalysisRetrievalsMap[key] = {};
                }
                state.workflowAnalysisRetrievalsMap[key][subKey] = true;
            } else {
                if (state.workflowAnalysisRetrievalsMap[key]) {
                    delete state.workflowAnalysisRetrievalsMap[key][subKey];

                    // Check if the main key has no inner keys left
                    if (Object.keys(state.workflowAnalysisRetrievalsMap[key]).length === 0) {
                        // If so, delete the main key from the map
                        delete state.workflowAnalysisRetrievalsMap[key];
                    }
                }
            }
        },
        setTaskMap: (state, action: PayloadAction<UpdateTaskCycleCountPayload>) => {
            const { key, count } = action.payload;
            const tmp = state.taskMap[key]
            tmp.cycleCount = count;
            state.taskMap[key] = tmp;
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
        removeAggregationFromWorkflowBuilderTaskMap: (state, action: PayloadAction<UpdateTaskMapPayload>) => {
            const { key, subKey, value } = action.payload;
                if (state.workflowBuilderTaskMap[key]) {
                    // Delete all subkeys from the value
                    Object.keys(state.workflowBuilderTaskMap[key]).forEach(subKey => {
                        delete state.workflowBuilderTaskMap[key][Number(subKey)];
                    });

                    delete state.workflowBuilderTaskMap[key];
                }
        },
    }
});

export const {
    setSearchContent,
    setGroupFilter,
    setUsernames,
    setAnalysisWorkflowInstructions,
    setAggregationWorkflowInstructions,
    setSearchResults,
    setWorkflows,
    setPlatformFilter,
    setAiTasks,
    setAddAnalysisView,
    setAddAggregationView,
    setAddRetrievalView,
    setAddAnalysisTasks,
    setAddAggregateTasks,
    setWorkflowBuilderTaskMap,
    removeAggregationFromWorkflowBuilderTaskMap,
    setTaskMap,
    setRetrievalName,
    setRetrievalGroup,
    setRetrievalPlatformGroups,
    setRetrievalKeywords,
    setRetrievalPlatform,
    setRetrievalUsernames,
    setRetrievalPrompt,
    setAddRetrievalTasks,
    setRetrievals,
    setAnalysisRetrievalsMap
} = aiSlice.actions;
export default aiSlice.reducer;
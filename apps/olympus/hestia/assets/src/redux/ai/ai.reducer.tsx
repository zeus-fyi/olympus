import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {
    Action,
    ActionMetric,
    ActionPlatformAccount,
    AiState,
    OrchestrationsAnalysis,
    PlatformSecretReference,
    Retrieval,
    SearchIndexerParams,
    TaskModelInstructions,
    UpdateTaskCycleCountPayload,
    UpdateTaskMapPayload
} from "./ai.types";

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
    retrievalsMap: {},
    retrieval: {
        retrievalName: '',
        retrievalGroup: '',
        retrievalKeywords: '',
        retrievalPlatform: '',
        retrievalUsernames: '',
        retrievalPrompt: '',
        retrievalPlatformGroups: '',
        discordFilters: {
            categoryName: '',
        }
    },
    retrievals: [],
    workflowAnalysisRetrievalsMap: {},
    workflowName: '',
    workflowGroupName: '',
    selectedWorkflows: [],
    runs: [],
    selectedRuns: [],
    searchIndexers: [],
    selectedSearchIndexers: [],
    searchIndexer: {
        searchID: 0,
        searchGroupName: '',
        platform: '',
        maxResults: 0,
        query: '',
        active: false,
    },
    platformSecretReference: {
        secretGroupName: 'mockingbird',
        secretKeyName: '',
    },
    selectedMainTab: 0,
    selectedMainTabBuilder: 0,
    action: {
        actionName: '',
        actionGroupName: '',
        actionType: '',
        actionStatus: '',
        actionPlatformAccounts: [],
        actionMetrics: [],
    },
    actions: [],
    actionMetric: {
        metricName: '',
        metricScoreThreshold: 1,
        metricPostActionMultiplier: 1,
    },
    actionPlatformAccount: {
        actionPlatformName: '',
        actionPlatformAccount: '',
    }
}

const aiSlice = createSlice({
    name: 'ai',
    initialState,
    reducers: {
        updateActionMetrics: (state, action: PayloadAction<ActionMetric[]>) => {
            state.action.actionMetrics = action.payload;
        },
        setActionPlatformAccount: (state, action: PayloadAction<ActionPlatformAccount>) => {
            state.actionPlatformAccount = action.payload;
        },
        setActionMetric: (state, action: PayloadAction<ActionMetric>) => {
            state.actionMetric = action.payload;
        },
        setAction: (state, action: PayloadAction<Action>) => {
            state.action = action.payload;
        },
        setActions: (state, action: PayloadAction<Action[]>) => {
            state.actions = action.payload;
        },
        setSelectedMainTab: (state, action: PayloadAction<number>) => {
            state.selectedMainTab = action.payload;
        },
        setSelectedMainTabBuilder: (state, action: PayloadAction<number>) => {
            state.selectedMainTabBuilder = action.payload;
        },
        setPlatformSecretReference: (state, action: PayloadAction<PlatformSecretReference>) => {
            state.platformSecretReference = action.payload;
        },
        setSearchIndexer: (state, action: PayloadAction<SearchIndexerParams>) => {
            state.searchIndexer = action.payload;
        },
        setSearchIndexers: (state, action: PayloadAction<SearchIndexerParams[]>) => {
            state.searchIndexers = action.payload;
        },
        setSelectedSearchIndexers: (state, action: PayloadAction<string[]>) => {
            state.selectedSearchIndexers = action.payload;
        },
        setWebRoutingGroup: (state, action: PayloadAction<string>) => {
            if (!state.retrieval.webFilters) {
                state.retrieval.webFilters = {
                    routingGroup: '',
                }
            }
            state.retrieval.webFilters.routingGroup = action.payload;
        },
        setRuns: (state, action: PayloadAction<OrchestrationsAnalysis[]>) => {
            state.runs = action.payload;
        },
        setSelectedRuns: (state, action: PayloadAction<string[]>) => {
            state.selectedRuns = action.payload;
        },
        setSelectedWorkflows: (state, action: PayloadAction<string[]>) => {
            state.selectedWorkflows = action.payload;
        },
        setWorkflowName: (state, action: PayloadAction<string>) => {
            state.workflowName = action.payload;
        },
        setWorkflowGroupName: (state, action: PayloadAction<string>) => {
            state.workflowGroupName = action.payload;
        },
        setRetrievalName: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalName = action.payload;
        },
        setRetrievalPlatformGroups: (state, action: PayloadAction<string>) => {
            state.retrieval.retrievalPlatformGroups = action.payload;
        },
        setDiscordOptionsCategoryName: (state, action: PayloadAction<string>) => {
            if (!state.retrieval.discordFilters) {
                state.retrieval.discordFilters = {
                    categoryName: '',
                }
            }
            state.retrieval.discordFilters.categoryName = action.payload;
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
        setRetrievals: (state, action: PayloadAction<Retrieval[]>) => {
            state.retrievals = action.payload;
        },
        setAddAnalysisTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAnalysisTasks = action.payload;
            for (let i = 0; i < state.addedAnalysisTasks.length; i++) {
                const task  = state.addedAnalysisTasks[i]
                if (task && task.taskID) {
                    if (task.cycleCount <= 0) {
                        task.cycleCount = 1;
                    }
                    state.taskMap[task.taskID] = task;
                }
            }
        },
        setAddAggregateTasks: (state, action: PayloadAction<TaskModelInstructions[]>) => {
            state.addedAggregateTasks = action.payload;
            for (let i = 0; i < state.addedAggregateTasks.length; i++) {
                const task  = state.addedAggregateTasks[i]
                if (task && task.taskID) {
                    if (task.cycleCount <= 0) {
                        task.cycleCount = 1;
                    }
                    state.taskMap[task.taskID] = task;
                }
            }
        },
        setAddRetrievalTasks: (state, action: PayloadAction<Retrieval[]>) => {
            state.addedRetrievals = action.payload;
            for (let i = 0; i < state.addedRetrievals.length; i++) {
                const retrieval  = state.addedRetrievals[i]
                if (retrieval && retrieval.retrievalID) {
                    state.retrievalsMap[retrieval.retrievalID] = retrieval;
                }
            }
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
            if (count <= 0) {
                tmp.cycleCount = 1;
            } else {
                tmp.cycleCount = count;
            }
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
    setAnalysisRetrievalsMap,
    setWorkflowGroupName,
    setWorkflowName,
    setSelectedWorkflows,
    setSelectedRuns,
    setRuns,
    setDiscordOptionsCategoryName,
    setWebRoutingGroup,
    setSelectedSearchIndexers,
    setSearchIndexers,
    setSearchIndexer,
    setPlatformSecretReference,
    setSelectedMainTab,
    setSelectedMainTabBuilder,
    setAction,
    setActions,
    setActionMetric,
    setActionPlatformAccount,
    updateActionMetrics,
} = aiSlice.actions;
export default aiSlice.reducer;
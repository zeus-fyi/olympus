export interface AiState {
    usernames: string;
    groupFilter: string;
    searchContentText: string;
    analysisWorkflowInstructions: string,
    aggregationWorkflowInstructions: string,
    searchResults: string;
    platformFilter: string;
    workflows: [];
    tasks: [];
    addAnalysisView: boolean;
    addAggregationView: boolean;
    addRetrievalView: boolean;
    addedAnalysisTasks: TaskModelInstructions[];
    addedAggregateTasks: TaskModelInstructions[];
    addedRetrievals: Retrieval[];
    workflowBuilderTaskMap: AggregateSubTasksMap
    taskMap: TaskMap;
    retrievalsMap: RetrievalsMap;
    retrieval: Retrieval;
    retrievals: Retrieval[];
    workflowAnalysisRetrievalsMap: AnalysisRetrievalsMap
}

export interface PostWorkflowsRequest {
    workflowName: string;
    stepSize: number;
    stepSizeUnit: string;
    models: TaskMap;
    aggregateSubTasksMap: AggregateSubTasksMap;
    analysisRetrievalsMap: AnalysisRetrievalsMap
}

export interface TaskModelInstructions {
    taskID?: number;
    group: string;
    model: string;
    taskType: string;
    taskGroup: string;
    taskName: string;
    maxTokens: number;
    tokenOverflowStrategy: string;
    prompt: string;
    cycleCount: number;
    retrievals?: AnalysisRetrievalsMap;
}

export interface TaskMap {
    [key: number]: TaskModelInstructions;
}

export interface AggregateSubTasksMap {
    [key: number]: { [innerKey: number]: boolean };
}

export type UpdateTaskMapPayload = {
    key: number;
    subKey: number;
    value: boolean;
};

export type UpdateTaskCycleCountPayload = {
    key: number;
    count: number;
};

export interface RetrievalsMap {
    [key: number]: Retrieval;
}

export interface AnalysisRetrievalsMap {
    [key: number]: { [innerKey: number]: boolean };
}

export interface Retrieval {
    retrievalID?: number;
    retrievalName: string;
    retrievalGroup: string;
    retrievalPrompt: string;
    retrievalKeywords: string;
    retrievalPlatform: string;
    retrievalUsernames: string;
    retrievalPlatformGroups: string;
}
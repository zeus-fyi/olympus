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
    addedAnalysisTasks: TaskModelInstructions[];
    addedAggregateTasks: TaskModelInstructions[];
    workflowBuilderTaskMap: AggregateSubTasksMap
    taskMap: TaskMap;
}

export interface PostWorkflowsRequest {
    workflowName: string;
    stepSize: number;
    stepSizeUnit: string;
    models: TaskMap;
    aggregateSubTasksMap: AggregateSubTasksMap;
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

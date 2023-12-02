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
}

export interface WorkflowModelInstructions {
    instructionType: string;
    model: string;
    maxTokens: number;
    tokenOverflowStrategy: string;
    cycleCount: number;
    prompt: string;
}

export interface PostWorkflowsRequest {
    workflowName: string;
    stepSize: number;
    stepSizeUnit: string;
    models: WorkflowModelInstructions[];
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
}

export interface AggregateSubTasksMap {
    [key: number]: { [innerKey: number]: boolean };
}

export type UpdateTaskMapPayload = {
    key: number;
    subKey: number;
    value: boolean;
};
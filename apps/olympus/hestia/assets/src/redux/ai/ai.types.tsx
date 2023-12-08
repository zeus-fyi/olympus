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
    workflowName: string;
    workflowGroupName: string;
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
    selectedWorkflows: string[];
    runs: OrchestrationsAnalysis[];
    selectedRuns: string[];
}

export interface PostWorkflowsRequest {
    workflowName: string;
    workflowGroupName: string;
    stepSize: number;
    stepSizeUnit: string;
    models: TaskMap;
    aggregateSubTasksMap?: AggregateSubTasksMap;
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

export interface AggregatedData {
    workflowResultId: number;
    responseId: number;
    sourceTaskId: number;
    taskName: string;
    taskType: string;
    runningCycleNumber: number;
    searchWindowUnixStart: number;
    searchWindowUnixEnd: number;
    model: string;
    prompt?: string; // or a more specific type if the structure of prompt is known
    metadata?: string; // or a more specific type if the structure of metadata is known
    completionChoices?: string; // similar to metadata, define a more specific type if possible
    promptTokens: number;
    completionTokens: number;
    totalTokens: number;
}

export interface OrchestrationsAnalysis {
    totalWorkflowTokenUsage: number;
    runCycles: number;
    aggregatedData: AggregatedData[];
    orchestrations: Orchestration;
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
    discordFilters?: DiscordFilters;
}

export interface DiscordFilters {
    categoryName: string;
}

export interface DeleteWorkflowsActionRequest {
    workflows: WorkflowTemplate[];
}

export interface PostWorkflowsActionRequest {
    action: string;
    unixStartTime: number;
    durationUnit: string;
    duration: number;
    customBasePeriod: boolean,
    customBasePeriodStepSize: number,
    customBasePeriodStepSizeUnit: string,
    workflows: WorkflowTemplate[];
}
export interface WorkflowTemplate {
    workflowID: number;
    workflowName: string;
    workflowGroup: string;
    fundamentalPeriod: number;
    fundamentalPeriodTimeUnit: string;
    tasks: Task[]; // Array of Task
}
export type Task = {
    taskName: string;
    taskType: string;
    model: string;
    prompt: string;
    cycleCount: number;
    retrievalName?: string;
    retrievalPlatform?: string;
};

export type Orchestration = {
    orchestrationID: number;
    active: boolean;
    groupName: string;
    type: string;
    instructions?: string;
    orchestrationName: string;
};

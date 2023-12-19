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
    selectedSearchIndexers: string[];
    runs: OrchestrationsAnalysis[];
    selectedRuns: string[];
    searchIndexers: SearchIndexerParams[];
    searchIndexer: SearchIndexerParams
    platformSecretReference: PlatformSecretReference;
    selectedMainTab: number;
    selectedMainTabBuilder: number;
    action: Action;
    actionPlatformAccount: ActionPlatformAccount
    actionMetric: ActionMetric;
    actions: Action[];
    evalFn: EvalFn
    evalFns: EvalFn[];
    evalMetric: EvalMetric;
    actionsEvalTrigger: EvalActionTrigger;
}

export interface EvalActionTrigger {
    evalState: string;
    evalCompletionStatus: string;
}

export interface EvalFn {
    evalID?: number;
    evalName: string;
    evalType: string;
    evalGroupName: string;
    evalModel?: string;
    evalFormat: string
    evalMetrics: EvalMetric[];
}
export interface EvalMetric {
    evalModelPrompt: string;
    evalMetricName: string;
    evalMetricResult: string;
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalMetricDataType: string;
    evalOperator: string;
    evalState: string;
}

export interface ActionPlatformAccount {
    actionPlatformName: string;
    actionPlatformAccount: string;
}

export interface Action {
    actionID?: number;
    actionName: string;
    actionGroupName: string;
    actionTriggerOn: string;
    actionType: string;
    actionStatus: string;
    actionEvals: EvalActionTrigger[];
    actionMetrics : ActionMetric[];
    actionPlatformAccounts: ActionPlatformAccount[];
}


export interface ActionMetric {
    metricName: string;
    metricScoreThreshold: number;
    metricPostActionMultiplier: number;
    metricOperator: string;
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
    responseFormat: string;
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
    orchestration: Orchestration;
}

export interface DeleteWorkflowsActionRequest {
    workflows: WorkflowTemplate[];
}

export interface PostRunsActionRequest {
    action: string;
    runs: Orchestration[];
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

export interface DiscordFilters {
    categoryTopic?: string;
    categoryName?: string;
    category?: string;
}

export interface AiSearchParams  {
    timeRange?: string;
    window?: Window;
    retrieval: Retrieval;
}

export interface Window {
    start?: Date;
    end?: Date;
    unixStartTime?: number;
    unixEndTime?: number;
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
    webFilters?: WebFilters;
}

export interface SearchIndexerParams {
    searchID: number;
    searchGroupName: string;
    maxResults: number;
    query: string;
    platform: string;
    active: boolean;
    discordOpts?: DiscordIndexerOpts;
}

export interface DiscordIndexerOpts {
    guildID: string;
    channelID: string;
}

export interface PlatformSecretReference {
    secretGroupName: string;
    secretKeyName: string;
}

export interface PostCreateOrUpdateSearchIndexerRequest {
    searchIndexer: SearchIndexerParams;
    platformSecretReference: PlatformSecretReference;
}

export interface PostSearchIndexerActionsRequest {
    action: string;
    searchIndexers: SearchIndexerParams[];
}

export interface TelegramIndexerOpts {
    SearchIndexerParams: SearchIndexerParams;
}

export interface WebFilters {
    routingGroup: string;
}
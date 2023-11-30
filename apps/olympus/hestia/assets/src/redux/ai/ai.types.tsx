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
    group: string;
    name: string;
    model: string;
    taskType: string;
    taskGroup: string;
    taskName: string;
    maxTokens: number;
    tokenOverflowStrategy: string;
    prompt: string;
}

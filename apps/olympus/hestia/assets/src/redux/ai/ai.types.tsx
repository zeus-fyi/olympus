export interface AiState {
    usernames: string;
    groupFilter: string;
    searchContentText: string;
    analysisWorkflowInstructions: string,
    aggregationWorkflowInstructions: string,
    searchResults: string;
    platformFilter: string;
    workflows: [{}];
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

import {EvalMetric} from "./ai.types.evals";
import {Orchestration} from "./ai.types";

export interface AggregatedData {
    workflowResultID: number;
    responseID: number;
    sourceTaskID: number;
    taskName: string;
    taskType: string;
    runningCycleNumber: number;
    chunkOffset: number;
    searchWindowUnixStart: number;
    searchWindowUnixEnd: number;
    iterationCount: number;
    skipAnalysis: boolean;
    model: string;
    prompt?: string; // or a more specific type if the structure of prompt is known
    metadata?: string; // or a more specific type if the structure of metadata is known
    completionChoices?: string; // similar to metadata, define a more specific type if possible
    promptTokens: number;
    completionTokens: number;
    totalTokens: number;
}

export interface RetrievalResult {
    workflowResultStrID: string;
    orchestrationID: number;
    retrievalID: number;
    retrievalName: string;
    runningCycleNumber: number;
    iterationCount: number;
    chunkOffset: number;
    searchWindowUnixStart: number;
    searchWindowUnixEnd: number;
    skipRetrieval: boolean;
    status: string;
    metadata?: any;
}


export interface OrchestrationsAnalysis {
    totalWorkflowTokenUsage: number;
    runCycles: number;
    progress: number;
    completeApiRequests: number;
    totalApiRequests: number;
    totalCsvCells: number;
    aggregatedData: AggregatedData[];
    orchestration: Orchestration;
    aggregatedEvalResults: EvalMetric[]; // Added array of EvalMetricsResult
    aggregatedRetrievalResults: RetrievalResult[]
}

export interface OrchDetailsMap {
    [index: string]: OrchestrationsAnalysis;
}

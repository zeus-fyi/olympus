// TypeScript interface for EvalMetricsResult
import {TriggerAction} from "./ai.types2";
import {JsonSchemaDefinition} from "./ai.types.schemas";

export interface EvalMetricsResult {
    evalName?: string;
    evalMetricName: string;
    evalMetricID?: number;
    evalMetricsResultId: number;
    evalMetricResult: string;
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalMetricDataType: string;
    evalOperator: string;
    evalState: string;
    runningCycleNumber: number;
    searchWindowUnixStart?: number;
    searchWindowUnixEnd?: number;
    evalResultOutcome: boolean;
    evalMetadata?: string; // Assuming json.RawMessage is defined elsewhere
}

export interface EvalFnMap {
    [key: number]: { [innerKey: number]: boolean };
}

export interface EvalMap {
    [key: number]: EvalFn;
}

export type UpdateEvalMapPayload = {
    evalID: number;
    evalTaskID: number;
    value: boolean;
};

export interface EvalFn {
    evalID?: number;
    evalTaskID?: number;
    evalName: string;
    evalType: string;
    evalGroupName: string;
    evalModel?: string;
    evalFormat: string
    evalCycleCount?: number;
    evalMetrics: EvalMetric[];
    triggerFunctions?: TriggerAction[];
    schemas: JsonSchemaDefinition[];
}

export interface EvalMetric {
    jsonSchemaID?: number;
    evalMetricID?: number;
    evalMetricResult: string;
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalOperator: string;
    evalState: string;
}
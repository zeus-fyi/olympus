// TypeScript interface for EvalMetricsResult
import {TriggerAction} from "./ai.types.retrievals";
import {JsonSchemaDefinition} from "./ai.types.schemas";

export interface EvalMetricsResult {
    evalMetricStrID?: string;
    evalMetricsResultStrID?: string;
    evalMetricsResult?: EvalMetricResult;
    evalMetricComparisonValues?: EvalMetricComparisonValues;
    evalOperator: string;
    evalState: string;
    runningCycleNumber: number;
    searchWindowUnixStart?: number;
    searchWindowUnixEnd?: number;
    evalResultOutcome: boolean;
    evalMetadata?: string; // Assuming json.RawMessage is defined elsewhere
}

export interface EvalMetricResult {
    evalMetricResultStrID?: string;
    evalResultOutcomeBool?: boolean;
    evalMetadata?: any; // Consider using a more specific type if the structure of evalMetadata is known
}

export interface EvalFnMap {
    [key: string]: { [innerKey: string]: boolean };
}

export interface EvalMap {
    [key: string]: EvalFn;
}

export type UpdateEvalMapPayload = {
    evalStrID: string;
    evalTaskStrID: string;
    value: boolean;
};

export interface EvalFn {
    evalStrID?: string;
    evalName: string;
    evalType: string;
    evalGroupName: string;
    evalModel: string;
    evalFormat: string;
    evalCycleCount?: number;
    evalMetrics: EvalMetric[];
    triggerFunctions?: TriggerAction[];
    schemas?: JsonSchemaDefinition[];
    schemasMap?: { [key: number]: JsonSchemaDefinition };
}

export interface EvalMetric {
    evalMetricStrID?: string;
    evalMetricResult?: EvalMetricResult;
    evalOperator: string;
    evalState: string;
    evalExpectedResultState: string;
    evalMetricComparisonValues?: EvalMetricComparisonValues;
}

export interface EvalMetricComparisonValues {
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalComparisonInteger?: number;
}


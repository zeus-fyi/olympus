// TypeScript interface for EvalMetricsResult
import {TriggerAction} from "./ai.types.retrievals";
import {JsonSchemaDefinition, JsonSchemaField} from "./ai.types.schemas";

export interface EvalMetric {
    evalMetricStrID?: string;
    evalName?: string;
    evalField?: JsonSchemaField;
    evalMetricResult?: EvalMetricResult;
    evalOperator: string;
    evalState: string;
    evalExpectedResultState: string;
    evalMetricComparisonValues?: EvalMetricComparisonValues;
}

export interface EvalMetricsResult {
    evalMetricsResultStrID?: string;
    evalResultOutcomeBool?: boolean;
    evalResultOutcomeStateStr?: string;
    evalMetadata?: any; // Consider using a more specific type if the structure of evalMetadata is known
    evalIterationCount?: number;
    runningCycleNumber?: number;
    searchWindowUnixStart?: number;
    searchWindowUnixEnd?: number;
}

export interface EvalMetricResult {
    evalResultOutcomeStateStr?: string;
    evalMetricResultStrID?: string;
    evalResultOutcomeBool?: boolean;
}

export interface EvalMetricComparisonValues {
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalComparisonInteger?: number;
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

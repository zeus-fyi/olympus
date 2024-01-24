// TypeScript interface for EvalMetricsResult
import {TriggerAction} from "./ai.types.retrievals";
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
    orgID?: number;
    userID?: number;
    evalName: string;
    evalType: string;
    evalGroupName: string;
    evalModel?: string;
    evalFormat: string;
    evalCycleCount?: number;
    evalMetrics: EvalMetric[];
    triggerFunctions?: TriggerAction[];
    schemas: JsonSchemaDefinition[];
    schemasMap: { [key: number]: JsonSchemaDefinition };
}

export interface EvalMetric {
    evalMetricID?: number;
    evalMetricResult?: EvalMetricResult;
    evalOperator: string;
    evalState: string;
    evalExpectedResultState: string;
    evalComparisonValues: EvalMetricComparisonValues;
}

export interface EvalMetricComparisonValues {
    evalComparisonBoolean?: boolean;
    evalComparisonNumber?: number;
    evalComparisonString?: string;
    evalComparisonInteger?: number;
}

export interface EvalMetricResult {
    evalMetricResultID?: number;
    evalResultOutcomeBool?: boolean;
    evalMetadata?: any; // Consider using a more specific type if the structure of evalMetadata is known
}

export interface FieldValue {
    intValue?: number;
    stringValue?: string;
    numberValue?: number;
    booleanValue?: boolean;
    intValueSlice?: number[];
    stringValueSlice?: string[];
    numberValueSlice?: number[];
    booleanValueSlice?: boolean[];
    isValidated?: boolean;
}
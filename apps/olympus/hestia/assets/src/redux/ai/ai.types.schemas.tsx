import {EvalMetric} from "./ai.types.evals";

export interface JsonSchemaDefinition {
    schemaID: number;
    schemaName: string;
    schemaGroup: string;
    schemaDescription: string;
    isObjArray: boolean;
    fields: JsonSchemaField[];
}

export interface JsonSchemaField {
    fieldID?: number;
    fieldName: string;
    fieldDescription: string;
    dataType: string;
    fieldValue?: FieldValue;
    evalMetrics?: EvalMetric[];
}

export interface AITaskJsonSchema {
    schemaID: number;
    taskID: number;
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
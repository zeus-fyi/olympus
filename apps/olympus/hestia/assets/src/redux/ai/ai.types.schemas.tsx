import {EvalMetric, FieldValue} from "./ai.types.evals";

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

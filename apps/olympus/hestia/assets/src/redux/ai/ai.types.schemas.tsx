
export interface JsonSchemaDefinition {
    schemaID: number;
    isObjArray: boolean;
    schemaName: string;
    schemaGroup: string;
    fields: JsonSchemaField[];
}

export interface JsonSchemaField  {
    fieldName: string;
    fieldDescription: string;
    dataType: string;
}

export interface AITaskJsonSchema {
    schemaID: number;
    taskID: number;
}

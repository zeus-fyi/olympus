export interface FlowState {
    uploadContentContacts: [];
    uploadContentTasks: [],
    csvHeaders: string[];
    promptHeaders: string[];
    results: [];
    stages: {};
    commandPrompts: {}
}

export const initialState: FlowState = {
    uploadContentContacts: [],
    uploadContentTasks: [],
    csvHeaders: [],
    promptHeaders: [],
    results: [],
    stages: {},
    commandPrompts: {}
}

export interface FlowAction {
    promptsCsv: [];
    contentContactsCsv: [];
    stages: {};
    contentContactsFieldMaps: {};
    commandPrompts: {}
}
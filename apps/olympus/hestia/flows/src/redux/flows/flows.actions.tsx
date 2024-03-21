export interface FlowState {
    uploadContentContacts: string;
    uploadContentTasks: string,
    csvHeaders: string[];
    promptHeaders: string[];
    results: [];
    stages: [];
}

export const initialState: FlowState = {
    uploadContentContacts: '',
    uploadContentTasks: '',
    csvHeaders: [],
    promptHeaders: [],
    results: [],
    stages: []
}

export interface FlowAction {
    promptsCsv: string;
    contentContactsCsv: string;
    stages: {};
    commandPrompts: {}
}
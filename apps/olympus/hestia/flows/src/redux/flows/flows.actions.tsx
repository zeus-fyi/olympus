export interface FlowState {
    uploadContentContacts: [];
    uploadContentTasks: [],
    csvHeaders: string[];
    promptHeaders: string[];
    results: [];
    stages: {
        linkedIn: boolean;
        googleSearch: boolean;
    };
    commandPrompts: {}
}

export const initialState: FlowState = {
    uploadContentContacts: [],
    uploadContentTasks: [],
    csvHeaders: [],
    promptHeaders: [],
    results: [],
    stages: {
        linkedIn: false,
        googleSearch: false
    },
    commandPrompts: {
        linkedin: '',
        googleSearch: ''
    }
}

// [key: string]: string;

export interface FlowAction {
    promptsCsv: [];
    contentContactsCsv: [];
    stages: {};
    contentContactsFieldMaps: {};
    commandPrompts: {}
}

export const googleSearchPromptOverride = {
    'biz-lead-google-search-summary': '', // use the settings for tasks
}

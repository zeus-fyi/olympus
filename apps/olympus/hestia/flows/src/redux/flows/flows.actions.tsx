export interface FlowState {
    uploadContentContacts: [];
    uploadContentTasks: [],
    csvHeaders: string[];
    promptHeaders: string[];
    results: [];
    stages: {
        linkedIn: boolean;
        linkedInBiz: boolean;
        googleSearch: boolean;
        validateEmails: boolean;
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
        linkedInBiz: false,
        googleSearch: false,
        validateEmails: false
    },
    commandPrompts: {
        linkedIn: 'Can you tell me their role and responsibilities?',
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

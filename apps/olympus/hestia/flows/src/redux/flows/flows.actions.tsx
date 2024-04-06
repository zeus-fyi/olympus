export interface FlowState {
    uploadContentContacts: [];
    promptsCsvContent: [],
    csvHeaders: string[];
    promptHeaders: string[];
    results: [];
    stages: {
        linkedIn: boolean;
        linkedInBiz: boolean;
        googleSearch: boolean;
        validateEmails: boolean;
        websiteScrape: boolean;
    };
    commandPrompts: {}
    previewCount: number
}

export const initialState: FlowState = {
    uploadContentContacts: [],
    promptsCsvContent: [],
    csvHeaders: [],
    promptHeaders: [],
    results: [],
    stages: {
        linkedIn: false,
        linkedInBiz: false,
        googleSearch: false,
        validateEmails: false,
        websiteScrape: false
    },
    commandPrompts: {
        linkedIn: '',
        linkedInBiz: '',
        googleSearch: '',
        websiteScrape: 'Can you tell me what the company does, and the industry they work in?'
    },
    previewCount: 0,
}

// [key: string]: string;

export interface FlowAction {
    previewCount: number;
    contentContactsCsv: [];
    contentContactsCsvStr: string;
    promptsCsv: [];
    promptsCsvStr: string;
    stages: {};
    contentContactsFieldMaps: {};
    commandPrompts: {}
}

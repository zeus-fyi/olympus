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
        googleSearch: '',
        websiteScrape: 'Can you tell me what the company does, and the industry they work in?'
    }
}

// [key: string]: string;

export interface FlowAction {
    contentContactsCsv: [];
    contentContactsCsvStr: string;
    promptsCsv: [];
    promptsCsvStr: string;
    stages: {};
    contentContactsFieldMaps: {};
    commandPrompts: {}
}

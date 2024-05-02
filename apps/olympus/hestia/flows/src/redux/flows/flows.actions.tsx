export interface FlowState {
    uploadContentContacts: [];
    promptsCsvContent: [],
    csvHeaders: string[];
    contactsCsvFilename: string;
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
    contactsCsvFilename: '',
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
    previewCount: 3,
}

// [key: string]: string;

export interface FlowAction {
    previewCount: number;
    contactsCsvFilename: string;
    contentContactsCsv: [];
    contentContactsCsvStr: string;
    promptsCsv: [];
    promptsCsvStr: string;
    stages: {};
    contentContactsFieldMaps: {};
    commandPrompts: {}
}

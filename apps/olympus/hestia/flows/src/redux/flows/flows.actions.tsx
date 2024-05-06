export interface FlowState {
    flowList: string[],
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
    stageColMap: StageColMap;        // Added new field
    stagePromptMap: StagePromptMap;  // Added new field
}

export const initialState: FlowState = {
    flowList: [],
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
    stageColMap: {},
    stagePromptMap: {},
    previewCount: 3,
}

export type UpdateTaskRelationshipPayload = {
    key: string;
    subKey: string;
};

export interface StageColMap {
    [key: string]: string ;
}

export interface StagePromptMap {
    [key: string]: string ;
}

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

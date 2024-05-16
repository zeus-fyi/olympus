export interface FlowState {
    adminFlowsMainTab: number;
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
    commandPrompts: {};
    customStages: {},
    previewCount: number,
    stageContactsOverrideMap: StageColMap;
    stageContactsMap: StageColMap;        // Added new field
    stagePromptMap: StagePromptMap;  // Added new field
}

export const initialState: FlowState = {
    adminFlowsMainTab: 0,
    flowList: [],
    uploadContentContacts: [],
    promptsCsvContent: [],
    csvHeaders: [],
    promptHeaders: [],
    results: [],
    contactsCsvFilename: '',
    customStages: {},
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
    stageContactsMap: {},
    stageContactsOverrideMap: {},
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
    commandPrompts: {};
    stagePromptMap: {};
    stageContactsMap:{};
    stageContactsOverrideMap: {};
}

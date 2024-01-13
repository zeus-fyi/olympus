
export interface Assistant {
    id: string;
    object: string;
    created_at: number | null;
    name: string;
    description:  string | null;
    model: string;
    instructions?: string;
    tools?: any;
    file_ids?: any[];
    metadata?: any;
}

export interface DiscordFilters {
    categoryTopic?: string;
    categoryName?: string;
    category?: string;
}

export interface RetrievalItemInstruction {
    retrievalPlatform: string;
    retrievalPlatformGroups?: string;
    retrievalKeywords?: string;
    retrievalUsernames?: string;
    retrievalPrompt?: string;
    discordFilters?: DiscordFilters;
    webFilters?: WebFilters;
    instructions?: string;
}

export interface WebFilters {
    routingGroup?: string;
    lbStrategy?: string;
}

export interface Retrieval {
    retrievalID?: number;
    retrievalName: string;
    retrievalGroup: string;
    retrievalItemInstruction: RetrievalItemInstruction;
}

export interface TriggerPlatformAccount {
    triggerPlatformName: string;
    triggerPlatformAccount: string;
}

export interface TriggerPlatform {
    triggerPlatformName: string;
    triggerPlatformAccount: TriggerPlatformAccount[];
}

export interface EvalActionTrigger {
    evalTriggerState: string;
    evalResultsTriggerOn: string;
}

export interface TriggerAction {
    triggerID?: number;
    triggerName: string;
    triggerGroup: string;
    triggerEnv: string;
    // triggerPlatforms: TriggerPlatform[];
    triggerActionsApprovals: TriggerActionsApproval[];
    evalTriggerActions: EvalActionTrigger[];
    evalTriggerAction: EvalActionTrigger;
    // actionMetrics : ActionMetric[];
    // actionPlatformAccounts: ActionPlatformAccount[];
}

export interface TriggerActionApprovalPutRequest {
    requestedState: string;
    triggerApproval: TriggerActionsApproval;
}
export type TriggerActionsApproval = {
    approvalID: number;
    evalID: number;
    triggerID: number;
    workflowResultID: number;
    approvalState: string;
    requestSummary: string;
    updatedAt: Date;
};

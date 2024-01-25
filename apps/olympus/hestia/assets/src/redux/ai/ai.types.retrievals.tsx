
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
    retrievalNegativeKeywords?: string;
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
    triggerStrID?: string;
    triggerName: string;
    triggerGroup: string;
    triggerAction: string;
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
    approvalStrID: number;
    evalStrID: number;
    triggerStrID: number;
    workflowResultStrID: number;
    approvalState: string;
    requestSummary: string;
    updatedAt: Date;
};

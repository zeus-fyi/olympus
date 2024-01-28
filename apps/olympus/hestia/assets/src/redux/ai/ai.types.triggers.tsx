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

import {Retrieval} from "./ai.types.retrievals";

export interface EvalActionTrigger {
    evalTriggerState: string;
    evalResultsTriggerOn: string;
}

export interface TriggerAction {
    triggerStrID?: string;
    triggerName: string;
    triggerGroup: string;
    triggerAction: string;
    triggerRetrievals: Retrieval[];
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

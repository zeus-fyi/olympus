
export interface Assistant {
    id: string;
    object: string;
    created_at: number | null;
    name: string;
    description:  string | null;
    model: string;
    instructions?: string;
    tools?: any;
    fileIDs?: any[];
    metadata?: any;
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
    evalTriggerActions: EvalActionTrigger[];
    // actionMetrics : ActionMetric[];
    // actionPlatformAccounts: ActionPlatformAccount[];
}

export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
    planUsageDetails: PlanUsageDetails;
    tableMetrics: any;
}

export interface Groups {
    [key: string]: string[];
}

export type PlanUsageDetails = {
    planName: string;
    computeUsage?: UsageMeter | null;
    tableUsage: TableUsage;
}

export type UsageMeter = {
    rateLimit: number;
    monthlyUsage: number;
    monthlyBudgetZU: number;
}

export type TableUsage = {
    endpointCount: number;
    tableCount: number;
    monthlyBudgetTableCount: number;
}

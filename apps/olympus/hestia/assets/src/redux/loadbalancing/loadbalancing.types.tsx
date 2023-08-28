
export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
    planUsageDetails: any;
    tableMetrics: any;
}

export interface Groups {
    [key: string]: string[];
}

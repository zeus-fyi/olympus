
export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
    planUsageDetails: any;
}

export interface Groups {
    [key: string]: string[];
}

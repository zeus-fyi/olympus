
export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
    userPlanInfo: any;
}

export interface Groups {
    [key: string]: string[];
}

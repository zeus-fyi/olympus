
export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
}

export interface Groups {
    [key: string]: string[];
}

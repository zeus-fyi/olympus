
export interface LoadBalancingState {
    routes: string[];
    groups: Groups;
    planUsageDetails: PlanUsageDetails;
    tableMetrics: TableMetricsSummary;
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
    currentRate : number;
    monthlyUsage: number;
    monthlyBudgetZU: number;
}

export type TableUsage = {
    endpointCount: number;
    tableCount: number;
    monthlyBudgetTableCount: number;
}

export interface TableMetricsSummary {
    tableName: string;
    routes: Z[];
    metrics: Record<string, TableMetric>;
}

export interface TableMetric {
    sampleCount: number;
    metricPercentiles: MetricSample[];
}

export interface MetricSample {
    percentile: number;
    latency: number;
}

export interface Z {
    Score: number;
    Member: any;
}


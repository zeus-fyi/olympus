import {ClusterConfigPage} from "./ClusterConfigPage";
import * as React from "react";
import {WorkloadConfigPage} from "./WorkloadConfigPage";

export const clusterBuilderSteps = [
    'Define Cluster',
    'Define Workloads',
];

export function wizardStepComponents(activeStep: number,
) {
    const steps = [
        <ClusterConfigPage />,
        <WorkloadConfigPage />
    ];
    return (steps[activeStep])
}
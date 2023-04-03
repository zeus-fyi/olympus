import {ClusterConfigPage} from "./ClusterConfigPage";
import * as React from "react";
import {WorkloadConfigPage} from "./WorkloadConfigPage";
import {WorkloadPreviewAndSubmitPage} from "./WorkloadPreviewAndSubmitPage";

export const clusterBuilderSteps = [
    'Define Cluster',
    'Define Workloads',
    'Preview & Create',
];

export function wizardStepComponents(activeStep: number,
) {

    const steps = [
        <ClusterConfigPage />,
        <WorkloadConfigPage />,
        <WorkloadPreviewAndSubmitPage />
    ];
    return (steps[activeStep])
}
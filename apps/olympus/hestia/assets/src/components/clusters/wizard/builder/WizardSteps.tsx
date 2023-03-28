import {DefineClusterClassParams} from "./DefineClusterClass";

export const clusterBuilderSteps = [
    'Define Cluster',
    'Define Component Base Workloads',
    'Define Skeleton Base Workloads',
    'Define Docker Image',
];

export function wizardStepComponents(activeStep: number,
) {
    const steps = [
        <DefineClusterClassParams />,
    ];
    return (steps[activeStep])
}
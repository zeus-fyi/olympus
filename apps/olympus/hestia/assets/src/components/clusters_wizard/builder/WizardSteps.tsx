import {DefineClusterClassParams} from "./DefineClusterClass";
import {DefineDockerParams} from "./DefineDockerImage";

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
        <DefineDockerParams />
    ];
    return (steps[activeStep])
}
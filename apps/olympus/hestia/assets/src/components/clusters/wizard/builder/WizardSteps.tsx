import {DefineClusterClassParams} from "./DefineClusterClass";
import {DefineDockerParams} from "./DefineDockerImage";
import {DefineClusterComponentBaseParams} from "./DefineComponentBases";

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
        <DefineClusterComponentBaseParams />,
        <DefineDockerParams />
    ];
    return (steps[activeStep])
}
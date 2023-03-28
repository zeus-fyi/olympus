import {DefineClusterClassParams} from "./DefineClusterClass";
import {DefineClusterComponentBaseParams} from "./DefineComponentBases";
import {AddSkeletonBaseDockerConfigs} from "./AddSkeletonBaseDockerConfigs";

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
        <AddSkeletonBaseDockerConfigs />
    ];
    return (steps[activeStep])
}
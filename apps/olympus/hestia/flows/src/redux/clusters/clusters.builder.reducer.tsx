import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {
    Cluster,
    ClusterPreview,
    ComponentBases,
    Container,
    DockerImage,
    Ingress,
    IngressPaths,
    Port,
    PVCTemplate,
    ResourceRequirements,
    SkeletonBase,
    SkeletonBases,
    VolumeMount
} from "./clusters.types";

interface ClusterBuilderState {
    clusterPreview: ClusterPreview;
    cluster: Cluster;
    selectedComponentBaseName: string;
    selectedSkeletonBaseName: string;
    selectedContainerName: string;
    selectedDockerImage: DockerImage;
    selectedClusterAppView: string;
    clusterViewEnabledToggle: boolean;
}

const initialState: ClusterBuilderState = {
    clusterPreview: {} as ClusterPreview,
    cluster: {
        clusterName: '',
        componentBases: {} as ComponentBases,
        ingressSettings: {authServerURL: 'aegis.zeus.fyi', host: 'host.zeus.fyi'} as Ingress,
        ingressPaths: {} as IngressPaths,
    },
    selectedComponentBaseName: '',
    selectedSkeletonBaseName: '',
    selectedContainerName: '',
    selectedDockerImage: {
        imageName: '',
        cmd: '',
        args: '',
        resourceRequirements: {cpu: '', memory: ''} as ResourceRequirements,
        ports: [{name: '', number: 0, protocol: 'TCP', ingressEnabledPort: false}] as Port[],
        volumeMounts: [{name: '', mountPath: ''}] as VolumeMount[]
    } as DockerImage,
    selectedClusterAppView: '',
    clusterViewEnabledToggle: false,
};

const clusterBuilderSlice = createSlice({
    name: 'clusterBuilder',
    initialState,
    reducers: {
        setClusterViewEnabledToggle: (state, action: PayloadAction<boolean>) => {
            state.clusterViewEnabledToggle = action.payload;
        },
        setSelectedClusterAppViewName: (state, action: PayloadAction<string>) => {
            state.selectedClusterAppView = action.payload;
        },
        setClusterPreview: (state, action: PayloadAction<ClusterPreview>) => {
            state.clusterPreview = action.payload;
        },
        setClusterName: (state, action: PayloadAction<string>) => {
            state.cluster.clusterName = action.payload;
        },
        setSelectedContainerName: (state, action: PayloadAction<string>) => {
            state.selectedContainerName = action.payload;
        },
        setSelectedComponentBaseName: (state, action: PayloadAction<string>) => {
            state.selectedComponentBaseName = action.payload;
        },
        setSelectedSkeletonBaseName: (state, action: PayloadAction<string>) => {
            state.selectedSkeletonBaseName = action.payload;
        },
        addComponentBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBases: SkeletonBases }>) => {
            const { componentBaseName, skeletonBases } = action.payload;
            state.cluster.componentBases[componentBaseName] = skeletonBases;
            state.cluster.ingressPaths[componentBaseName] = {path: '/', pathType: 'ImplementationSpecific'};
        },
        removeComponentBase: (state, action: PayloadAction<string>) => {
            const key = action.payload;
            if (state.cluster.componentBases[key]) {
                delete state.cluster.componentBases[key];
            } else {
                console.error(`Component base not found: ${key}`);
            }
            if (state.cluster.ingressPaths[key]) {
                delete state.cluster.ingressPaths[key];
            }
        },
        addSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; skeletonBase: SkeletonBase }>) => {
            const { componentBaseName, skeletonBaseName, skeletonBase } = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName] = skeletonBase;
        },
        setConfigMapKey: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; key: string; value: string; }>) => {
            const {componentBaseName, skeletonBaseName, key, value} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].configMap[key] = value;
        },
        removeConfigMapKey: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; key: string;}>) => {
            const {componentBaseName, skeletonBaseName, key} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            delete state.cluster.componentBases[componentBaseName][skeletonBaseName].configMap[key]
        },
        setIngressHost: (state, action: PayloadAction<{ host: string }>) => {
            const {host} = action.payload;
            state.cluster.ingressSettings.host = host;
        },
        setIngressPath: (state, action: PayloadAction<{ componentBaseName: string; path: string }>) => {
            const {componentBaseName,  path} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.ingressPaths[componentBaseName].path = path;
        },
        setIngressPathType: (state, action: PayloadAction<{ componentBaseName: string; pathType: string }>) => {
            const {componentBaseName,  pathType} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.ingressPaths[componentBaseName].pathType = pathType;
        },
        removeIngressPath: (state, action: PayloadAction<{ componentBaseName: string}>) => {
            const {componentBaseName} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            delete state.cluster.ingressPaths[componentBaseName]
        },
        setIngressAuthServerURL: (state, action: PayloadAction<{ authServerURL: string }>) => {
            const { authServerURL} = action.payload;
            state.cluster.ingressSettings.authServerURL = authServerURL;
        },
        setStatefulSetReplicaCount: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; replicaCount: number }>) => {
            const {componentBaseName, skeletonBaseName, replicaCount} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].statefulSet.replicaCount = replicaCount;
        },
        setStatefulSetPVC: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; pvcIndex: number, pvc: PVCTemplate }>) => {
            const { componentBaseName, skeletonBaseName, pvcIndex, pvc } = action.payload;
            const skeletonBase = state.cluster.componentBases[componentBaseName]?.[skeletonBaseName];
            if (!skeletonBase) {
                console.error(`Skeleton base not found: ${skeletonBaseName}`);
                return;
            }
            const statefulSet = skeletonBase.statefulSet;
            if (!statefulSet) {
                console.error(`Stateful set not found in skeleton base: ${skeletonBaseName}`);
                return;
            }
            if (pvcIndex < 0 || pvcIndex >= statefulSet.pvcTemplates.length) {
                console.error(`Invalid pvc index: ${pvcIndex}`);
                return;
            }
            if (pvcIndex >= 0) {
                statefulSet.pvcTemplates[pvcIndex] = pvc;
            } else {
                statefulSet.pvcTemplates.push(pvc);
            }
        },
        addStatefulSetPVC: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; pvc: PVCTemplate }>) => {
            const { componentBaseName, skeletonBaseName, pvc } = action.payload;
            const skeletonBase = state.cluster.componentBases[componentBaseName]?.[skeletonBaseName];
            if (!skeletonBase) {
                console.error(`Skeleton base not found: ${skeletonBaseName}`);
                return;
            }
            const statefulSet = skeletonBase.statefulSet;
            if (!statefulSet) {
                console.error(`Stateful set not found in skeleton base: ${skeletonBaseName}`);
                return;
            }
            statefulSet.pvcTemplates.push(pvc);
        },
        removeStatefulSetPVC: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; pvcIndex: number, pvc: PVCTemplate}>) => {
            const { componentBaseName, skeletonBaseName, pvcIndex, pvc} = action.payload;
            const skeletonBase = state.cluster.componentBases[componentBaseName]?.[skeletonBaseName];
            if (!skeletonBase) {
                console.error(`Skeleton base not found: ${skeletonBaseName}`);
                return;
            }
            const statefulSet = skeletonBase.statefulSet;
            if (!statefulSet) {
                console.error(`Stateful set not found in skeleton base: ${skeletonBaseName}`);
                return;
            }

            if (pvcIndex !== undefined && (pvcIndex < 0 || pvcIndex >= skeletonBase.statefulSet.pvcTemplates.length)) {
                console.error(`Invalid pvc index: ${pvcIndex}`);
                return;
            }
            if (pvcIndex !== undefined) {
                // Remove port at specified index
                skeletonBase.statefulSet.pvcTemplates.splice(pvcIndex, 1);
            } else if (pvc !== undefined) {
                // Add or update port
                if (pvcIndex !== undefined) {
                    // Update existing port
                    skeletonBase.statefulSet.pvcTemplates[pvcIndex] = pvc;
                } else {
                    // Add new port
                    skeletonBase.statefulSet.pvcTemplates.push(pvc);
                }
            }
        },
        setDeploymentReplicaCount: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; replicaCount: number }>) => {
            const {componentBaseName, skeletonBaseName, replicaCount} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].deployment.replicaCount = replicaCount;
        },
        toggleStatefulSetWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addStatefulSet: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addStatefulSet} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addStatefulSet = addStatefulSet;
            if (state.cluster.componentBases[componentBaseName][skeletonBaseName].addDeployment && addStatefulSet) {
                state.cluster.componentBases[componentBaseName][skeletonBaseName].addDeployment = false;
                state.cluster.componentBases[componentBaseName][skeletonBaseName].deployment.replicaCount = 0;
                state.cluster.componentBases[componentBaseName][skeletonBaseName].containers = {};
            }
        },
        toggleDeploymentWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addDeployment: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addDeployment} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addDeployment = addDeployment;
            if (state.cluster.componentBases[componentBaseName][skeletonBaseName].addStatefulSet && addDeployment) {
                state.cluster.componentBases[componentBaseName][skeletonBaseName].addStatefulSet = false;
                state.cluster.componentBases[componentBaseName][skeletonBaseName].statefulSet.replicaCount = 0;
                state.cluster.componentBases[componentBaseName][skeletonBaseName].containers = {};
            }
        },
        toggleServiceWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addService: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addService} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addService = addService;
        },
        toggleIngressWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addIngress: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addIngress} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addIngress = addIngress;
        },
        toggleConfigMapWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addConfigMap: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addConfigMap} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addConfigMap = addConfigMap;
        },
        toggleAddServiceMonitorWorkloadSelectionOnSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string; addServiceMonitor: boolean }>) => {
            const {componentBaseName, skeletonBaseName, addServiceMonitor} = action.payload;
            if (!state.cluster.componentBases[componentBaseName]) {
                state.cluster.componentBases[componentBaseName] = {};
            }
            state.cluster.componentBases[componentBaseName][skeletonBaseName].addServiceMonitor = addServiceMonitor;
        },
        removeSkeletonBase: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string }>) => {
            const { componentBaseName, skeletonBaseName } = action.payload;
            if (state.cluster.componentBases[componentBaseName][skeletonBaseName]) {
                delete state.cluster.componentBases[componentBaseName][skeletonBaseName];
            } else {
                console.error(`Skeleton base not found: ${skeletonBaseName}`);
            }
        },
        addContainer: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; container: Container }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, container } = action.payload;
            if (!state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]) {
                console.error(`SkeletonBase not found: ${skeletonBaseKey}`);
                return;
            }
            state.cluster.componentBases[componentBaseKey][skeletonBaseKey].containers[containerName] = container;
        },
        setContainerInit: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; isInitContainer: boolean }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, isInitContainer } = action.payload;
            if (!state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]) {
                console.error(`SkeletonBase not found: ${skeletonBaseKey}`);
                return;
            }
            state.cluster.componentBases[componentBaseKey][skeletonBaseKey].containers[containerName].isInitContainer = isInitContainer;
        },
        removeContainer: (state, action: PayloadAction<{ componentBaseName: string; skeletonBaseName: string, containerName: string}>) => {
            const { componentBaseName, skeletonBaseName, containerName} = action.payload;
            if (state.cluster.componentBases[componentBaseName][skeletonBaseName].containers[containerName]) {
                delete state.cluster.componentBases[componentBaseName][skeletonBaseName].containers[containerName];
            } else {
                console.error(`Container not found: ${containerName}`);
            }
        },
        setDockerImageCmd: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; cmd: string }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, cmd } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found in container setDockerImageCmd: ${containerName}`);
                return;
            }
            dockerImage.cmd = cmd
        },
        setDockerImage: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; dockerImageKey: string;}>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey } = action.payload;
            const container = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName];
            if (!container) {
                console.error(`Docker image not found in container setDockerImage: ${containerName}`);
                return;
            }
            container.dockerImage.imageName = dockerImageKey
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.imageName = dockerImageKey
        },
        setSelectedDockerImage: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string;}>) => {
            const { componentBaseKey, skeletonBaseKey, containerName} = action.payload;
            const container = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName];
            if (!container) {
                console.error(`Docker image not found in container: ${containerName}`);
                return;
            }
            state.selectedDockerImage = container.dockerImage
        },
        setDockerImageCmdArgs: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; args: string}>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, args } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found in container: ${containerName}`);
                return;
            }
            dockerImage.args = args
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.args = args
        },
        setDockerImageCpuResourceRequirement: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; cpu: string}>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, cpu } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found in container: ${containerName}`);
                return;
            }
            dockerImage.resourceRequirements.cpu = cpu
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.resourceRequirements.cpu = cpu
        },
        setDockerImageMemoryResourceRequirement: (state, action: PayloadAction<{ componentBaseKey: string; skeletonBaseKey: string; containerName: string; memory: string}>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, memory } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found in container: ${containerName}`);
                return;
            }
            dockerImage.resourceRequirements.memory = memory
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.resourceRequirements.memory = memory
        },
        setDockerImagePort: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            portIndex: number;
            port: Port;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, portIndex, port } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            if (portIndex < 0 || portIndex >= dockerImage.ports.length) {
                console.error(`Invalid port index: ${portIndex}`);
                return;
            }
            const updatedPorts = dockerImage.ports.map((p, i) => {
                if (i === portIndex) {
                    return port;
                } else if (port.ingressEnabledPort && p.ingressEnabledPort) {
                    return { ...p, ingressEnabledPort: false };
                } else {
                    return p;
                }
            });
            dockerImage.ports = updatedPorts;
            const selectedDockerImage = state.selectedDockerImage;
            selectedDockerImage.ports = updatedPorts;
        },
        addDockerImagePort: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            port: Port;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, port } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            // if (!port.name || port.number === 0) {
            //     console.error(`Invalid port: ${port}`);
            //     return;
            // }
            dockerImage.ports.push(port);
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.ports.push(port);
        },
        removeDockerImagePort: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            portIndex?: number;
            port?: Port;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, portIndex, port } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            const selectedDockerImage = state.selectedDockerImage

            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            if (portIndex !== undefined && (portIndex < 0 || portIndex >= dockerImage.ports.length)) {
                console.error(`Invalid port index: ${portIndex}`);
                return;
            }
            if (portIndex !== undefined && port === undefined) {
                // Remove port at specified index
                dockerImage.ports.splice(portIndex, 1);
                selectedDockerImage.ports.splice(portIndex, 1);
            } else if (port !== undefined) {
                // Add or update port
                if (portIndex !== undefined) {
                    // Update existing port
                    dockerImage.ports[portIndex] = port;
                    selectedDockerImage.ports[portIndex] = port;
                } else {
                    // Add new port
                    dockerImage.ports.push(port);
                    selectedDockerImage.ports.push(port);
                }
            }
        },
        setDockerImageVolumeMount: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            volumeMountIndex: number;
            volumeMount: VolumeMount;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, volumeMountIndex, volumeMount } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            if (volumeMountIndex < 0 || volumeMountIndex >= dockerImage.volumeMounts.length) {
                console.error(`Invalid volume mount index: ${volumeMountIndex}`);
                return;
            }
            dockerImage.volumeMounts[volumeMountIndex] = volumeMount;
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.volumeMounts[volumeMountIndex] = volumeMount;
        },
        addDockerImageVolumeMount: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            volumeMount: VolumeMount;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, volumeMount } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }
            dockerImage.volumeMounts.push(volumeMount);
            const selectedDockerImage = state.selectedDockerImage
            selectedDockerImage.volumeMounts.push(volumeMount);
        },
        removeDockerImageVolumeMount: (state, action: PayloadAction<{
            componentBaseKey: string;
            skeletonBaseKey: string;
            containerName: string;
            dockerImageKey: string;
            volumeMountIndex?: number;
            volumeMount?: VolumeMount;
        }>) => {
            const { componentBaseKey, skeletonBaseKey, containerName, dockerImageKey, volumeMountIndex, volumeMount } = action.payload;
            const dockerImage = state.cluster.componentBases[componentBaseKey]?.[skeletonBaseKey]?.containers[containerName].dockerImage;
            const selectedDockerImage = state.selectedDockerImage;

            if (!dockerImage) {
                console.error(`Docker image not found: ${dockerImageKey}`);
                return;
            }

            if (volumeMountIndex !== undefined && (volumeMountIndex < 0 || volumeMountIndex >= dockerImage.volumeMounts.length)) {
                console.error(`Invalid volume mount index: ${volumeMountIndex}`);
                return;
            }

            if (volumeMountIndex !== undefined && volumeMount === undefined) {
                // Remove volume mount at specified index
                dockerImage.volumeMounts.splice(volumeMountIndex, 1);
                selectedDockerImage.volumeMounts.splice(volumeMountIndex, 1);
            } else if (volumeMount !== undefined) {
                // Add or update volume mount
                if (volumeMountIndex !== undefined) {
                    // Update existing volume mount
                    dockerImage.volumeMounts[volumeMountIndex] = volumeMount;
                    selectedDockerImage.volumeMounts[volumeMountIndex] = volumeMount;
                } else {
                    // Add new volume mount
                    dockerImage.volumeMounts.push(volumeMount);
                    selectedDockerImage.volumeMounts.push(volumeMount);
                }
            }
        },
    },
});

export const {
    setClusterName,
    addComponentBase,
    removeComponentBase,
    addSkeletonBase,
    setSelectedContainerName,
    removeSkeletonBase,
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName,
    addContainer,
    setDockerImagePort,
    setDockerImageCmd,
    removeContainer,
    setDockerImage,
    setDockerImageCmdArgs,
    setSelectedDockerImage,
    removeDockerImagePort,
    addDockerImagePort,
    setDockerImageVolumeMount,
    addDockerImageVolumeMount,
    removeDockerImageVolumeMount,
    setContainerInit,
    setDockerImageCpuResourceRequirement,
    setDockerImageMemoryResourceRequirement,
    toggleStatefulSetWorkloadSelectionOnSkeletonBase,
    toggleDeploymentWorkloadSelectionOnSkeletonBase,
    toggleServiceWorkloadSelectionOnSkeletonBase,
    toggleIngressWorkloadSelectionOnSkeletonBase,
    toggleConfigMapWorkloadSelectionOnSkeletonBase,
    toggleAddServiceMonitorWorkloadSelectionOnSkeletonBase,
    removeStatefulSetPVC,
    setStatefulSetPVC,
    setStatefulSetReplicaCount,
    setDeploymentReplicaCount,
    addStatefulSetPVC,
    setConfigMapKey,
    removeConfigMapKey,
    setIngressHost,
    setIngressPath,
    setIngressPathType,
    removeIngressPath,
    setClusterPreview,
    setIngressAuthServerURL,
    setSelectedClusterAppViewName,
    setClusterViewEnabledToggle
} = clusterBuilderSlice.actions;

export default clusterBuilderSlice.reducer;
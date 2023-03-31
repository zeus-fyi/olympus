import * as React from "react";
import {Box, Button, Card, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import {AddSkeletonBaseDockerConfigs, SelectedSkeletonBaseName} from "./AddSkeletonBaseDockerConfigs";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    toggleAddServiceMonitorWorkloadSelectionOnSkeletonBase,
    toggleConfigMapWorkloadSelectionOnSkeletonBase,
    toggleDeploymentWorkloadSelectionOnSkeletonBase,
    toggleIngressWorkloadSelectionOnSkeletonBase,
    toggleServiceWorkloadSelectionOnSkeletonBase,
    toggleStatefulSetWorkloadSelectionOnSkeletonBase,
} from "../../../../redux/clusters/clusters.builder.reducer";

export function WorkloadConfigPage(props: any) {
    const {} = props;
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedComponentBase = cluster.componentBases?.[selectedComponentBaseName]?.[selectedSkeletonBaseName] ?? '';
    const addDeployment = selectedComponentBase?.addDeployment
    const addStatefulSet = selectedComponentBase?.addStatefulSet
    const addService = selectedComponentBase?.addService
    const addIngress = selectedComponentBase?.addIngress
    const addServiceMonitor = selectedComponentBase?.addServiceMonitor
    const addConfigMap = selectedComponentBase?.addConfigMap

    const onToggleStatefulSet = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addStatefulSet: !addStatefulSet,
        }
        dispatch(toggleStatefulSetWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleDeployment = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addDeployment: !addDeployment,
        }
        dispatch(toggleDeploymentWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleService = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addService: !addService,
        }
        dispatch(toggleServiceWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleConfigMap = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addConfigMap: !addConfigMap,
        }
        dispatch(toggleConfigMapWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleServiceMonitor = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addServiceMonitor: !addServiceMonitor,
        }
        dispatch(toggleAddServiceMonitorWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleIngress = () => {
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addIngress: !addIngress,
        }
        dispatch(toggleIngressWorkloadSelectionOnSkeletonBase(refObj));
    };
    return (
        <div>
            <Stack direction="row" spacing={2}>
                <Card sx={{ maxWidth: 500 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Workload Config
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Sets Infra and App Configs
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="column" spacing={2}>
                        <Box mt={2}>
                            <SelectedComponentBaseName />
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseName />
                        </Box>
                        {!addDeployment && (
                            <Button variant="contained" onClick={onToggleDeployment}>
                                Add Deployment
                            </Button>
                        )}
                        {addDeployment && (
                            <Stack direction="row" spacing={2}>
                            <Button variant="contained" color="primary">
                                    Deployment
                                </Button>
                                <Button variant="contained" onClick={onToggleDeployment}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {!addStatefulSet && (
                            <Button variant="contained" onClick={onToggleStatefulSet}>
                                Add StatefulSet
                            </Button>
                        )}
                        {addStatefulSet && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary">
                                    StatefulSet
                                </Button>
                                <Button variant="contained" onClick={onToggleStatefulSet}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {!addService && (
                            <Button variant="contained" onClick={onToggleService}>
                                Add Service
                            </Button>
                        )}
                        {addService && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary">
                                    Service
                                </Button>
                                <Button variant="contained" onClick={onToggleService}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {!addIngress && (
                            <Button variant="contained" onClick={onToggleIngress}>
                                Add Ingress
                            </Button>
                        )}
                        {addIngress && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary">
                                    Ingress
                                </Button>
                                <Button variant="contained" onClick={onToggleIngress}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {!addServiceMonitor && (
                            <Button variant="contained" onClick={onToggleServiceMonitor}>
                                Add ServiceMonitor
                            </Button>
                        )}
                        {addServiceMonitor && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary">
                                    ServiceMonitor
                                </Button>
                                <Button variant="contained" onClick={onToggleServiceMonitor}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        </Stack>
                    </Container>
                </Card>
                <AddSkeletonBaseDockerConfigs />
            </Stack>
        </div>
    );
}
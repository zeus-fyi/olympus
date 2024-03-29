import * as React from "react";
import {useState} from "react";
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
import {ServiceView} from "./ServiceView";
import {IngressView} from "./IngressView";
import {ConfigMapView} from "./ConfigMapView";

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
    const [viewField, setViewField] = useState('');

    const onClickView = (fieldName: string) => {
        setViewField(fieldName)
    };

    const onToggleStatefulSet = () => {
        if (!addStatefulSet) {
            setViewField('statefulSet')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addStatefulSet: !addStatefulSet,
        }
        dispatch(toggleStatefulSetWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleDeployment = () => {
        if (!addDeployment) {
            setViewField('deployment')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addDeployment: !addDeployment,
        }
        dispatch(toggleDeploymentWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleService = () => {
        if (!addService) {
            setViewField('service')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addService: !addService,
        }
        dispatch(toggleServiceWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleConfigMap = () => {
        if (!addConfigMap) {
            setViewField('configMap')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addConfigMap: !addConfigMap,
        }
        dispatch(toggleConfigMapWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleServiceMonitor = () => {
        if (!addServiceMonitor) {
            setViewField('serviceMonitor')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addServiceMonitor: !addServiceMonitor,
        }
        dispatch(toggleAddServiceMonitorWorkloadSelectionOnSkeletonBase(refObj));
    };
    const onToggleIngress = () => {
        if (!addIngress) {
            setViewField('ingress')
        } else {
            setViewField('')
        }
        const refObj = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            addIngress: !addIngress,
        }
        dispatch(toggleIngressWorkloadSelectionOnSkeletonBase(refObj));
    };

    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }
    let show = Object.keys(cluster.componentBases).length > 0
    return (
        <div>
            {show && Object.keys(cluster.componentBases?.[selectedComponentBaseName]).length > 0 && (
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
                    <Stack direction="column" spacing={2}>
                        <Container maxWidth="xl" sx={{}}>
                        <Box mt={2}>
                            <SelectedComponentBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        </Container>
                    </Stack>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="column" spacing={2}>
                        {!addDeployment && (
                            <Button variant="contained" onClick={onToggleDeployment}>
                                Add Deployment
                            </Button>
                        )}
                        {addDeployment && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary" onClick={() => onClickView('deployment')}>
                                    View Deployment
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
                                <Button variant="contained" color="primary" onClick={() => onClickView('statefulSet')}>
                                    View StatefulSet
                                </Button>
                                <Button variant="contained" onClick={onToggleStatefulSet}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {!addConfigMap && (
                            <Button variant="contained" onClick={onToggleConfigMap}>
                                Add ConfigMap
                            </Button>
                        )}
                        {addConfigMap && (
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" color="primary" onClick={() => onClickView('configMap')}>
                                    View ConfigMap
                                </Button>
                                <Button variant="contained" onClick={onToggleConfigMap}>
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
                                <Button variant="contained" color="primary" onClick={() => onClickView('service')}>
                                    View Service
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
                                <Button variant="contained" color="primary" onClick={() => onClickView('ingress')}>
                                    View Ingress
                                </Button>
                                <Button variant="contained" onClick={onToggleIngress}>
                                    Remove
                                </Button>
                            </Stack>
                        )}
                        {/*{!addServiceMonitor && (*/}
                        {/*    <Button variant="contained" onClick={onToggleServiceMonitor}>*/}
                        {/*        Add ServiceMonitor*/}
                        {/*    </Button>*/}
                        {/*)}*/}
                        {/*{addServiceMonitor && (*/}
                        {/*    <Stack direction="row" spacing={2}>*/}
                        {/*        <Button variant="contained" color="primary" onClick={() => onClickView('serviceMonitor')}>*/}
                        {/*            View ServiceMonitor*/}
                        {/*        </Button>*/}
                        {/*        <Button variant="contained" onClick={onToggleServiceMonitor}>*/}
                        {/*            Remove*/}
                        {/*        </Button>*/}
                        {/*    </Stack>*/}
                        {/*)}*/}
                        </Stack>
                    </Container>
                </Card>
                {(viewField === 'deployment' || viewField === 'statefulSet') && (
                    <div>
                        <AddSkeletonBaseDockerConfigs viewField={viewField} />
                    </div>
                )}
                {viewField === 'configMap' && (
                    <div>
                        <ConfigMapView />
                    </div>
                )}
                {viewField === 'service' && (
                    <div>
                        <ServiceView addStatefulSet={addStatefulSet} addDeployment={addStatefulSet} />
                    </div>
                )}
                {viewField === 'ingress' && (
                    <div>
                        <IngressView />
                    </div>
                )}
                {/*{viewField === 'serviceMonitor' && (*/}
                {/*    <div>*/}
                {/*        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>*/}
                {/*        </Container>*/}
                {/*    </div>*/}
                {/*)}*/}
            </Stack>
            )}
        </div>
    );
}


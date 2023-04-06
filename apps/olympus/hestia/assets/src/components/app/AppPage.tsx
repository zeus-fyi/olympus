import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import * as React from "react";
import {useEffect, useState} from "react";
import {
    Box,
    Button,
    Card,
    CardContent,
    CircularProgress,
    Container,
    Divider,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack
} from "@mui/material";
import {clustersApiGateway} from "../../gateway/clusters";
import {ClusterPreview} from "../../redux/clusters/clusters.types";
import {
    setClusterPreview,
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../redux/apps/apps.reducer";
import Typography from "@mui/material/Typography";
import YamlTextFieldAppPage from "./YamlFormattedTextAppPage";

export function AppPage(props: any) {
    const {} = props;
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const [viewField, setViewField] = useState('');
    const [previewType, setPreviewType] = useState('');
    const clusterPreview = useSelector((state: RootState) => state.apps.clusterPreview);
    const selectedComponentBaseName = useSelector((state: RootState) => state.apps.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.apps.selectedSkeletonBaseName);
    const [addDeployment, setAddDeployment] = useState(false);
    const [addConfigMap, setAddConfigMap] = useState(false);
    const [addIngress, setAddIngress] = useState(false);
    const [addService, setAddService] = useState(false);
    const [addStatefulSet, setAddStatefulSet] = useState(false);

    useEffect(() => {
        const skeletonBasePreview = clusterPreview?.componentBases?.[selectedComponentBaseName]?.[selectedSkeletonBaseName];
        if (skeletonBasePreview) {
            setAddDeployment(skeletonBasePreview.deployment !== null);
            setAddConfigMap(skeletonBasePreview.configMap !== null);
            setAddIngress(skeletonBasePreview.ingress !== null);
            setAddService(skeletonBasePreview.service !== null);
            setAddStatefulSet(skeletonBasePreview.statefulSet !== null);
        }
    }, [clusterPreview, selectedComponentBaseName, selectedSkeletonBaseName]);

    const [configError, setConfigError] = useState('');

    let buttonLabelCreate;
    let buttonDisabledCreate;
    let statusMessageCreate;
    const [requestCreateStatus, setCreateRequestStatus] = useState('');
    switch (requestCreateStatus) {
        case 'pending':
            buttonLabelCreate = <CircularProgress size={20} />;
            buttonDisabledCreate = true;
            break;
        case 'success':
            buttonLabelCreate = 'Register Cluster';
            buttonDisabledCreate = true;
            statusMessageCreate = 'Cluster definition generated successfully!';
            break;
        case 'error':
            buttonLabelCreate = 'Resubmit';
            buttonDisabledCreate = false;
            statusMessageCreate = 'An error occurred while submitting, there\'s likely a problem with your configuration, check that your ports, resource values, etc are valid. ' + configError;
            break;
        default:
            buttonLabelCreate = 'Register Cluster';
            buttonDisabledCreate = true;
            break;
    }

    const onClickView = (newPreviewType: string) => {
        setPreviewType(newPreviewType);
    }
    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }

    const onClickCreate = async () => {
        try {
            setCreateRequestStatus('pending');
            let res: any = await clustersApiGateway.createCluster(cluster)
            const cp =  res.data as ClusterPreview;
            const statusCode = res.status;
            if (statusCode === 200 || statusCode === 204) {
                dispatch(setClusterPreview(cp));
                setCreateRequestStatus('success');
            } else {
                setConfigError('')
                setCreateRequestStatus('error');
            }
        } catch (e) {
            setCreateRequestStatus('error');
        }
    }

    return (
        <div>
            <Stack direction="row" spacing={2}>
                <Card sx={{ minWidth: 250, maxWidth: 300 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            App Config
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mb: 4 }}>
                        <Box mt={2}>
                            <SelectedComponentBaseNameAppPage onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseNameAppsPage onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                    </Container>
                    <Divider/>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="column" spacing={2}>
                            {addDeployment && (
                                <Stack direction="row" spacing={2}>
                                    <Button variant="contained" color="primary" onClick={() => onClickView('deployment')}>
                                        View Deployment
                                    </Button>
                                </Stack>
                            )}
                            {addStatefulSet && (
                                <Stack direction="row" spacing={2}>
                                    <Button variant="contained" color="primary" onClick={() => onClickView('statefulSet')}>
                                        View StatefulSet
                                    </Button>
                                </Stack>
                            )}
                            {addConfigMap && (
                                <Stack direction="row" spacing={2}>
                                    <Button variant="contained" color="primary" onClick={() => onClickView('configMap')}>
                                        View ConfigMap
                                    </Button>
                                </Stack>
                            )}
                            {addService && (
                                <Stack direction="row" spacing={2}>
                                    <Button variant="contained" color="primary" onClick={() => onClickView('service')}>
                                        View Service
                                    </Button>
                                </Stack>
                            )}
                            {addIngress && (
                                <Stack direction="row" spacing={2}>
                                    <Button variant="contained" color="primary" onClick={() => onClickView('ingress')}>
                                        View Ingress
                                    </Button>
                                </Stack>
                            )}
                        </Stack>
                    </Container>
                    <Divider/>
                    <Container maxWidth="xl" sx={{ mb: 4 }}>
                        <Box mt={2}>
                            <Button variant="contained" onClick={onClickCreate} disabled={buttonDisabledCreate}>
                                {buttonLabelCreate}
                            </Button>
                            {statusMessageCreate && (
                                <Typography variant="body2" color={requestCreateStatus === 'error' ? 'error' : 'success'}>
                                    {statusMessageCreate}
                                </Typography>
                            )}
                        </Box>
                    </Container>
                </Card>
                <YamlTextFieldAppPage previewType={previewType}/>
            </Stack>
        </div>
    );
}


export function SelectedComponentBaseNameAppPage(props: any) {
    const {onChangeComponentOrSkeletonBase} = props;
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.apps.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.apps.selectedComponentBaseName);
    const onAccessComponentBase = (selectedComponentBaseName: string) => {
        dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
        const skeletonBaseName = Object.keys(cluster.componentBases[selectedComponentBaseName])[0];
        dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
        // Add a check to see if the `containers` field exists

        onChangeComponentOrSkeletonBase();
    };

    let show = Object.keys(cluster.componentBases).length > 0;
    return (
        <div>
            {show &&
                <FormControl sx={{mb: 1}} variant="outlined" style={{ minWidth: '100%' }}>
                    <InputLabel id="network-label">Cluster Bases</InputLabel>
                    <Select
                        labelId="componentBase-label"
                        id="componentBase"
                        value={selectedComponentBaseName}
                        label="Component Base"
                        onChange={(event) => onAccessComponentBase(event.target.value as string)}
                        sx={{ width: '100%' }}
                    >
                        {Object.keys(cluster.componentBases).map((key: any, i: number) => (
                            <MenuItem key={i} value={key}>
                                {key}
                            </MenuItem>))
                        }
                    </Select>
                </FormControl>
            }
        </div>);
}

export function SelectedSkeletonBaseNameAppsPage(props: any) {
    const { onChangeComponentOrSkeletonBase}  = props;
    const dispatch = useDispatch();
    const skeletonBaseName = useSelector((state: RootState) => state.apps.selectedSkeletonBaseName);
    const componentBaseName = useSelector((state: RootState) => state.apps.selectedComponentBaseName);
    const cluster = useSelector((state: RootState) => state.apps.cluster);

    useEffect(() => {
        dispatch(setSelectedComponentBaseName(componentBaseName));
        dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
    }, [dispatch,skeletonBaseName, componentBaseName]);

    const onAccessSkeletonBase = (selectedSkeletonBaseName: string) => {
        dispatch(setSelectedSkeletonBaseName(selectedSkeletonBaseName));
        onChangeComponentOrSkeletonBase();
    };

    if (cluster.componentBases === undefined) {
        return <div></div>
    }

    const skeletonBaseKeys = cluster.componentBases[componentBaseName];
    const show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }
    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Workload Bases</InputLabel>
            <Select
                labelId="skeletonBase-label"
                id="skeletonBase"
                value={skeletonBaseName}
                label="Skeleton Base"
                onChange={(event) => onAccessSkeletonBase(event.target.value as string)}
                sx={{ width: '100%' }}
            >
                {show && Object.keys(skeletonBaseKeys).map((key: any, i: number) => (
                    <MenuItem key={i} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}
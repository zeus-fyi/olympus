import {Box, Button, Card, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import * as React from "react";
import {useEffect, useState} from "react";
import {SelectedSkeletonBaseName} from "./AddSkeletonBaseDockerConfigs";
import YamlTextField from "./YamlFormattedTextPage";
import {clustersApiGateway} from "../../../../gateway/clusters";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {ClusterPreview} from "../../../../redux/clusters/clusters.types";
import {setClusterPreview} from "../../../../redux/clusters/clusters.builder.reducer";

export function WorkloadPreviewAndSubmitPage(props: any) {
    const {} = props;
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const [viewField, setViewField] = useState('');
    const [previewType, setPreviewType] = useState('');
    const clusterPreview = useSelector((state: RootState) => state.clusterBuilder.clusterPreview);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
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


    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');
    const dispatch = useDispatch();

    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Logged in successfully';
            buttonDisabled = true;
            statusMessage = 'Logged in successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while generating preview, please try again. If you continue having issues please email alex@zeus.fyi';
            break;
        default:
            buttonLabel = 'Login';
            buttonDisabled = false;
            break;
    }
    const onClickView = (newPreviewType: string) => {
        setPreviewType(newPreviewType);
    }
    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }

    const onClickPreviewCreate = async () => {
        try {
            setRequestStatus('pending');
            let res: any = await clustersApiGateway.previewCreateCluster(cluster)
            const cp =  res.data as ClusterPreview;
            const statusCode = res.status;
            if (statusCode === 200 || statusCode === 204) {
                dispatch(setClusterPreview(cp));
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
            }
        } catch (e) {
            setRequestStatus('error');
        }
    }

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
                    <Container maxWidth="xl" sx={{ mb: 4 }}>
                        <Box mt={2}>
                            <SelectedComponentBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                    </Container>
                    <Container maxWidth="xl" sx={{ mb: 4 }}>
                        <Box mt={2}>
                            <Button variant="contained" onClick={onClickPreviewCreate}>
                                Generate Preview
                            </Button>
                        </Box>
                        {/*<Box mt={2}>*/}
                        {/*    <Button variant="contained">*/}
                        {/*        Create Cluster*/}
                        {/*    </Button>*/}
                        {/*</Box>*/}
                    </Container>
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
                </Card>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <YamlTextField previewType={previewType}/>
                </Container>
            </Stack>
        </div>
    );
}
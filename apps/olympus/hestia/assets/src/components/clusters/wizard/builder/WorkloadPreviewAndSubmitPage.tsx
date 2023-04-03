import {Box, Button, Card, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import * as React from "react";
import {useState} from "react";
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
    const clusterPreview = useSelector((state: RootState) => state.clusterBuilder.clusterPreview);
    const selectedComponentBase =  useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName =  useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const clusterPreviewComponentBases = clusterPreview.componentBases;
    let selectedContent = '';

    if (clusterPreviewComponentBases !== undefined && Object.keys(clusterPreviewComponentBases).length > 0) {
        if (clusterPreviewComponentBases[selectedSkeletonBaseName] !== undefined && Object.keys(clusterPreviewComponentBases[selectedSkeletonBaseName]).length > 0) {
            //selectedContent = clusterPreview.componentBases[selectedComponentBase]?.[selectedSkeletonBaseName]?.statefulSet ?? '';
        }
    }
    const [viewField, setViewField] = useState('');
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
    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }

    const onClickPreviewCreate = async () => {
        console.log(cluster)
        try {
            setRequestStatus('pending');
            let res: any = await clustersApiGateway.previewCreateCluster(cluster)
            const cp =  res.data as ClusterPreview;
            console.log(cp)
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
                        <Box mt={2}>
                            <Button variant="contained">
                                Create Cluster
                            </Button>
                        </Box>
                    </Container>
                </Card>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <YamlTextField selectedContent={selectedContent}/>
                </Container>
            </Stack>
        </div>
    );
}
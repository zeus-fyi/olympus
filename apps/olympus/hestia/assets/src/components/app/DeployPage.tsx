import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import {
    Box,
    Card,
    CardActions,
    CardContent,
    CircularProgress,
    createTheme,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack
} from "@mui/material";
import {useState} from "react";
import {appsApiGateway} from "../../gateway/apps";
import {useParams} from "react-router-dom";
import {ThemeProvider} from "@mui/material/styles";
import {ResourceRequirementsTable} from "./ResourceRequirementsTable";

const mdTheme = createTheme();

export function DeployPage(props: any) {
    const [cloudProvider, setCloudProvider] = useState('do');
    const [region, setRegion] = useState('nyc1');
    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const params = useParams();

    const [requestStatus, setRequestStatus] = useState('');

    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Deploy More';
            buttonDisabled = false;
            statusMessage = 'Deployment in Progress';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while attempting to deploy.';
            break;
        default:
            buttonLabel = 'Deploy';
            buttonDisabled = true;
            break;
    }
    const handleDeploy = async () => {
        try {
            setRequestStatus('pending');
            const response = await appsApiGateway.deployApp(params.id as string, {});
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};

    function handleChangeSelectPathType(cloudProvider: string) {
        setCloudProvider(cloudProvider);
    }

    return (
        <div>
            <ThemeProvider theme={mdTheme}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Card sx={{ maxWidth: 700 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Deployment & Management
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Currently, you cannot deploy clusters without getting authorization manually first, until we have automated billing setup.
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="row" spacing={2}>
                            <FormControl sx={{ mr: 1 }} fullWidth variant="outlined">
                                <InputLabel key={`cloudProviderLabel`} id={`cloudProvider`}>
                                    Cloud Provider
                                </InputLabel>
                                <Select
                                    labelId={`cloudProviderLabel`}
                                    id={`cloudProvider`}
                                    name="cloudProvider"
                                    value={cloudProvider}
                                    onChange={(event) => handleChangeSelectPathType(event.target.value)}
                                    label="Cloud Provider"
                                >
                                    <MenuItem value="do">DigitalOcean</MenuItem>
                                </Select>
                            </FormControl>
                            <FormControl sx={{ mr: 1 }} fullWidth variant="outlined">
                                <InputLabel key={`regionLabel`} id={`region`}>
                                    Region
                                </InputLabel>
                                <Select
                                    labelId={`regionLabel`}
                                    id={`region`}
                                    name="region"
                                    value={region}
                                    onChange={(event) => handleChangeSelectPathType(event.target.value)}
                                    label="Region"
                                >
                                    <MenuItem value="nyc1">Nyc1</MenuItem>
                                </Select>
                            </FormControl>
                            <CardActions >
                                <Button variant="contained" onClick={handleDeploy} disabled={buttonDisabled}>{buttonLabel}</Button>
                                {statusMessage && (
                                    <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                                        {statusMessage}
                                    </Typography>
                                )}
                            </CardActions>
                        </Stack>
                    </Container>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Config Options
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Select which components and config options you want to deploy for this app.
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
                        <Box sx={{ mt: 2, display: 'flex' }}>
                            <ResourceRequirementsTable />
                        </Box>
                    </Container>
                </Card>
            </Container>
            {/*<Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>*/}
            {/*    <ResourceRequirementsTable />*/}
            {/*</Container>*/}
            {/*<Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>*/}
            {/*    <Card sx={{ maxWidth: 700 }}>*/}
            {/*        <CardContent>*/}
            {/*            <Typography gutterBottom variant="h5" component="div">*/}
            {/*                App Resource Requirements*/}
            {/*            </Typography>*/}
            {/*            <Typography variant="body2" color="text.secondary">*/}
            {/*                TODO*/}
            {/*            </Typography>*/}
            {/*        </CardContent>*/}
            {/*    </Card>*/}
            {/*</Container>*/}
            </ThemeProvider>

        </div>
    );
}
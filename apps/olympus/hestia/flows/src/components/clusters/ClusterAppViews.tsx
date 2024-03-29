import {
    Box,
    CardContent,
    Divider,
    FormControl,
    FormControlLabel,
    FormGroup,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    Switch
} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Button from "@mui/material/Button";
import {clustersApiGateway} from "../../gateway/clusters";
import {useDispatch} from "react-redux";
import {setClusterViewEnabledToggle} from "../../redux/clusters/clusters.builder.reducer";

export const ClusterViews = (props: any) => {
    const { pageView, setPageView, appName, setAppName, clusters, allClusters } = props;
    const [statusMessage, setStatusMessage] = React.useState("");
    const [appLabelingEnabled, setAppLabelingEnabled] = React.useState(true);
    const dispatch = useDispatch();

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPageView(event.target.checked);
        dispatch(setClusterViewEnabledToggle(event.target.checked));
    };

    const handleLabelToggle = (event: React.ChangeEvent<HTMLInputElement>) => {
        setAppLabelingEnabled(event.target.checked);
    };

    const uniqueAppNames: string[] = Array.from(new Set(allClusters.map((cluster: { clusterClassName: string }) => cluster.clusterClassName)));
    function handleChangeSelectAppView(appName: string) {
        if (appName === "-all") {
            setAppName("");
            return;
        }
        setAppLabelingEnabled(true)
        setAppName(appName);
    }

    const onClickRolloutUpgrade = async (clusterClassName: string) => {
        if (appName === "-all") {
            return;
        }
        try {
            const response = await clustersApiGateway.deployUpdateFleet(clusterClassName, appLabelingEnabled);
            const statusCode = response.status;
            if (statusCode === 202) {
                setStatusMessage(`Cluster fleet ${clusterClassName} update in progress`);
            } else if (statusCode === 200){
                setStatusMessage(`Cluster fleet ${clusterClassName} already up to date`);
            } else {
                setStatusMessage(`Cluster fleet ${clusterClassName} had an unexpected response: status code ${statusCode}`);
            }
        } catch (e) {
            setStatusMessage(`Cluster fleet ${clusterClassName} failed to update`);
        }
    }

    const onClickRolloutRestart = async (clusterClassName: string) => {
        if (appName === "-all") {
            return;
        }
        try {
            const response = await clustersApiGateway.deployRolloutRestartFleet(clusterClassName);
            const statusCode = response.status;
            if (statusCode === 202) {
                setStatusMessage(`Cluster fleet ${clusterClassName} rollout restart in progress`);
            } else if (statusCode === 200){
                setStatusMessage(`Cluster fleet ${clusterClassName} rollout restart already in progress`);
            } else {
                setStatusMessage(`Cluster fleet ${clusterClassName} had an unexpected response: status code ${statusCode}`);
            }
        } catch (e) {
            setStatusMessage(`Cluster fleet ${clusterClassName} failed to update`);
        }
    }

    return (
        <div>
            <Stack direction={"row"} spacing={2} alignItems={"center"}>
                {pageView ? (
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Cluster Apps View
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            This view shows all the apps that are deployed. If you upgrade the fleet it will by
                            default set app toleration tainting, so make sure you have nodes that can schedule
                            this app if you didn't deploy them from the UI otherwise disable it with the toggle.
                        </Typography>
                    </CardContent>
                ) : (
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Cluster View
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            This view shows all the cloud clusters you have and can interact with.
                        </Typography>
                    </CardContent>
                )}
                <FormGroup>
                    <FormControlLabel control={
                        <Switch
                        checked={pageView}
                        onChange={handleChange}
                        color="primary"
                        name="pageView"
                        inputProps={{ 'aria-label': 'toggle page view' }}
                    />} label={pageView ? 'Apps' : 'All'} />
                </FormGroup>
            </Stack>
            {pageView ? (
                <div>
                <Box mr={2} ml={2} mt={2} mb={4}>
                    <Stack direction={"row"} spacing={2} alignItems={"center"}>
                    <FormControl sx={{  }} fullWidth variant="outlined">
                        <InputLabel key={`appNameLabel`} id={`appName`}>
                            App View
                        </InputLabel>
                        <Select
                            labelId={`appNameLabel`}
                            id={`appName`}
                            name="appName"
                            value={appName}
                            onChange={(event) => handleChangeSelectAppView(event.target.value)}
                            label="App View"
                        >
                            <MenuItem key={'all'} value={'-all'}>{"all"}</MenuItem>
                            {uniqueAppNames.map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                        </Select>
                    </FormControl>
                    </Stack>
                </Box>
                    {appName !== '' && appName !== '-all' && (
                        <Divider />
                    )}
                    <div>
                        {appName !== '' && appName !== '-all' && (
                            <div>
                                <Stack direction={"row"} mr={2} ml={2} mt={2} mb={2} alignItems={"center"}>
                                        <Button  fullWidth onClick={() => onClickRolloutUpgrade(appName)} variant="contained">
                                            Upgrade Fleet
                                        </Button>
                                         <Box mr={2}></Box>
                                        <Button fullWidth onClick={() => onClickRolloutRestart(appName)} variant="contained">
                                            Rollout Restart Fleet
                                        </Button>
                               </Stack>
                            </div>
                        )}
                        <Box ml={2}>
                            <div>{statusMessage}</div>
                        </Box>
                    </div>
                    <Box mr={2} ml={2} mt={2} mb={2}>
                        <div>
                            {appName !== '' && appName !== '-all' && (
                                <FormGroup>
                                    <FormControlLabel control={
                                        <Switch
                                            checked={appLabelingEnabled}
                                            onChange={handleLabelToggle}
                                            color="primary"
                                            name="appLabelingEnabled"
                                            inputProps={{ 'aria-label': 'toggle page view' }}
                                        />} label={appLabelingEnabled ? 'App Tainting Enabled' : 'App Tainting Disabled'} />
                                </FormGroup>
                            )}
                        </div>
                    </Box>
                </div>

            ) : (<div></div>)}
        </div>
    );
};
import {Box, CardContent, FormControl, InputLabel, MenuItem, Select, Stack, Switch} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Button from "@mui/material/Button";
import {clustersApiGateway} from "../../gateway/clusters";

export const ClusterViews = (props: any) => {
    const { pageView, setPageView, appName, setAppName, clusters, allClusters } = props;
    const [statusMessage, setStatusMessage] = React.useState("");
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPageView(event.target.checked);
    };
    const uniqueAppNames: string[] = Array.from(new Set(allClusters.map((cluster: { clusterClassName: string }) => cluster.clusterClassName)));
    function handleChangeSelectAppView(appName: string) {
        if (appName === "-all") {
            setAppName("");
            return;
        }
        setAppName(appName);
    }

    const onClickRolloutUpgrade = async (clusterClassName: string) => {
        if (appName === "-all") {
            return;
        }
        try {
            const response = await clustersApiGateway.deployUpdateFleet(clusterClassName);
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
    return (
        <div>
            <Stack direction={"row"} spacing={2} alignItems={"center"}>
                {pageView ? (
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Cluster Apps View
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            This view shows all the apps that are deployed
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
                <Switch
                    checked={pageView}
                    onChange={handleChange}
                    color="primary"
                    name="pageView"
                    inputProps={{ 'aria-label': 'toggle page view' }}
                />
            </Stack>
            {pageView ? (
                <Box mr={2} ml={2} mt={2} mb={2}>
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
                    <div>
                        {appName !== '' && appName !== '-all' && (
                            <Button onClick={() => onClickRolloutUpgrade(appName)} variant="contained">
                                Upgrade Fleet
                            </Button>
                        )}
                        <div>{statusMessage}</div>
                    </div>
                    </Stack>
                </Box>
            ) : (<div></div>)}
        </div>
    );
};
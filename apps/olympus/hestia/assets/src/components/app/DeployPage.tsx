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
    Divider,
    FormControl,
    IconButton,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    TextField
} from "@mui/material";
import * as React from "react";
import {useEffect, useState} from "react";
import {appsApiGateway} from "../../gateway/apps";
import {useParams} from "react-router-dom";
import {ThemeProvider} from "@mui/material/styles";
import {createDiskResourceRequirements, ResourceRequirementsTable} from "./ResourceRequirementsTable";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {Nodes} from "../../redux/apps/apps.types";
import {Add, Remove} from "@mui/icons-material";
import {
    setCloudProvider,
    setCluster,
    setClusterPreview,
    setDeployServersCount,
    setNodes,
    setRegion,
    setSelectedComponentBaseName,
    setSelectedDisk,
    setSelectedNode,
    setSelectedSkeletonBaseName
} from "../../redux/apps/apps.reducer";
import {AppConfigsTable} from "./AppConfigTable";

const mdTheme = createTheme();

interface NodeMap {
    [resourceID: number]: Nodes;
}

export function DeployPage(props: any) {
    const {app} = props
    const cloudProviderRegionsResourcesMap = useSelector((state: RootState) => state.apps.cloudRegionResourceMap);
    const cloudProvider = useSelector((state: RootState) => state.apps.selectedCloudProvider);
    const region = useSelector((state: RootState) => state.apps.selectedRegion);
    const node = useSelector((state: RootState) => state.apps.selectedNode);
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const resourceRequirements = createDiskResourceRequirements(cluster);
    let nodes = useSelector((state: RootState) => state.apps.nodes);
    const nodeMap: NodeMap = {};

    const count =  useSelector((state: RootState) => state.apps.deployServersCount);
    const [freeTrial, setFreeTrial] = useState(false);
    const disk =  useSelector((state: RootState) => state.apps.selectedDisk);
    const dispatch = useDispatch();
    const params = useParams();

    useEffect(() => {
        async function fetchData() {
            let cluster;
            let clusterPreview;
            let selectedComponentBaseName;
            let selectedSkeletonBaseName;
            try {
                let id = params.id as string;
                if (app === "avax") {
                    id = "avax"
                }
                if (app === "ethereumEphemeralBeacons") {
                    id = "ethereumEphemeralBeacons"
                }
                if (app === "microservice") {
                    id = "microservice"
                }
                if (app === "sui") {
                    id = "sui"
                }
                const response = await appsApiGateway.getPrivateAppDetails(id);
                clusterPreview = await response.clusterPreview;
                dispatch(setClusterPreview(clusterPreview));
                cluster = await response.cluster;
                dispatch(setCluster(cluster));
                const cBases = await response.cluster.componentBases
                const cb = Object.keys(cBases)
                if (cb.length > 0) {
                    selectedComponentBaseName = cb[0];
                    dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
                    const sbs = Object.keys(response.cluster.componentBases[selectedComponentBaseName])
                    if (sbs.length > 0) {
                        selectedSkeletonBaseName = sbs[0];
                        dispatch(setSelectedSkeletonBaseName(selectedSkeletonBaseName));
                    }
                }
                if (response.cloudRegionResourceMap) {
                    dispatch(setNodes( cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.nodes || []))
                    const cloudProviderKeys = Object.keys(cloudProviderRegionsResourcesMap);
                    const firstCloudProvider = cloudProviderKeys.length > 0 ? cloudProviderKeys[0] : '';
                    dispatch(setCloudProvider(firstCloudProvider));

                    const regionKeys = Object.keys(cloudProviderKeys);
                    const firstRegion = regionKeys.length > 0 ? regionKeys[0] : '';

                    dispatch(setRegion(firstRegion));
                }
                nodes = cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.nodes || [];
                const filteredNodes = cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.nodes || [];
                filteredNodes.forEach((node) => {
                    if (node.resourceID === 0) {
                        return;
                    }
                    nodeMap[node.resourceID] = node;
                });
                return response;
            } catch (e) {
            }
        }
        // if (params.id) {
        //     if (getPrefix(params.id.toString()) === 'sui') {
        //         if (getSuffix(params.id.toString()) === 'aws') {
        //             dispatch(setRegion('us-west-1');
        //             setCloudProvider('aws');
        //         }
        //         if (getSuffix(params.id.toString()) === 'gcp') {
        //             setRegion('us-central1');
        //             setCloudProvider('gcp');
        //         }
        //         if (getSuffix(params.id.toString()) === 'do') {
        //             setRegion('nyc1');
        //             setCloudProvider('do');
        //         }
        //     }
        // }

    }, [params.id, nodes, nodeMap, cloudProviderRegionsResourcesMap]);

    // Dynamically fetch the regions for the selected cloud provider
    const regions = cloudProviderRegionsResourcesMap[cloudProvider] || {};
    const nodesInRegion = cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.nodes || [];
    // const isNodeInMap = (resourceID) => nodesInRegion.some(node => node.resourceID === resourceID);
    function isNodeInMap(resourceID: number) {
        return nodesInRegion.some(node => node.resourceID === resourceID);
    }
    // Generate the MenuItem components for each region
    const regionMenuItems = Object.keys(regions).map(region => (
        <MenuItem key={region} value={region}>{region}</MenuItem>
    ));

    const handleDiskTypeChange = (event: any) => {
        const disksInRegion = cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.disks || [];
        const selectedDisk = disksInRegion.find(disk => `${disk.type}-${disk.subType}` === event.target.value);
        if (selectedDisk) {
           dispatch(setSelectedDisk(selectedDisk));
        }
    };
    const disksInRegion = cloudProviderRegionsResourcesMap[cloudProvider]?.[region]?.disks || [];
    const diskCloudRegionMenuItems = disksInRegion.map(disk => (
        <MenuItem key={disk.type+'-'+disk.subType} value={disk.type+'-'+disk.subType}>{disk.type+'-'+disk.subType}</MenuItem>
    ));

    const handleIncrement = () => {
        dispatch(setDeployServersCount(count + 1));
    };

    const handleDecrement = () => {
        if (count - 1 < 0) {
            dispatch(setDeployServersCount(0))
            return;
        }
        dispatch(setDeployServersCount(count - 1));
    };

    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');
    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Refresh To Deploy More';
            buttonDisabled = true;
            statusMessage = 'Deployment in Progress';
            break;
        case 'missingBilling':
            buttonLabel = 'Try Free For One Hour';
            buttonDisabled = count === 0 && totalCost() <= 500;
            statusMessage = 'No payment methods have been set. You can set a payment option on the billing page.\n You can deploy it for free for one hour, but if a payment method isn\'t set it will automatically delete after one hour. For free trials the total monthly cost must be <= $500'
            break;
        case 'outOfCredits':
            buttonLabel = 'Try Free For One Hour';
            buttonDisabled = count === 0 && totalCost() <= 500;
            statusMessage = 'No payment methods have been set. You can set a payment option on the billing page.\n You can only deploy a maximum of one app at any time using the free trial. For free trials the total monthly cost must be <= $500. You\'ve reached the maximum credits,' +
                ' if you need more time or a higher free trial limit email alex@zeus.fyi'
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = count === 0;
            statusMessage = 'An error occurred while attempting to deploy.';
            break;
        default:

            if (cluster && cluster.clusterName === 'sui' ) {
                buttonLabel = 'Select Config';
                buttonDisabled = true;
            } else {
                buttonLabel = 'Deploy';
                buttonDisabled = false;
            }
            break;
    }
    const handleDeploy = async () => {
        try {
            setRequestStatus('pending');
            const namespaceAlias = cluster.clusterName;
            const payload = {
                "cloudCtxNs": {
                    "cloudProvider": cloudProvider,
                    "region": region,
                },
                "nodes": node,
                "count": count,
                "namespaceAlias": namespaceAlias,
                "cluster": cluster,
                "resourceRequirements": resourceRequirements,
                "freeTrial": freeTrial,
                "monthlyCost": totalCost()
            }
            const response = await appsApiGateway.deployApp(payload);
            if (response.status === 200 || response.status === 202 || response.status === 204) {
                setRequestStatus('success');
                return
            } else if (response.status === 403) {
                setRequestStatus('missingBilling');
                setFreeTrial(true)
                return
            } else if (response.status === 412) {
                setRequestStatus('outOfCredits');
                setFreeTrial(true)
                return
            } else {
                setRequestStatus('error');
                return
            }
        } catch (error: any) {
            setRequestStatus('error');
            const status: number = error.response.status;
            if (status === 403) {
                setRequestStatus('missingBilling');
                setFreeTrial(true)
            } else if (status === 412) {
                setRequestStatus('outOfCredits');
                // Disable the button for 30 seconds
                setFreeTrial(true)
            } else {
                setRequestStatus('error');
            }
        }};

    function handleChangeSelectRegion(region: string) {
        const regionsMap = cloudProviderRegionsResourcesMap[cloudProvider];
        const disks = regionsMap[region]?.disks || [];
        // Check if there are any disks available and set related states
        if (disks.length > 0) {
            const firstDisk = disks[0];

            dispatch(setSelectedDisk(firstDisk));
        } else {
            // Handle the case where no disks are available for the first region of the selected provider
            // This could be setting to default values or handling as a special case
        }
        dispatch(setRegion(region));
    }
    function handleChangeSelectCloudProvider(cloudProvider: string) {
        dispatch(setCloudProvider(cloudProvider));
        const regionsMap = cloudProviderRegionsResourcesMap[cloudProvider];

        if (regionsMap && Object.keys(regionsMap).length > 0) {
            // Get the first region's key from the regions map
            const firstRegion = Object.keys(regionsMap)[0];
            dispatch(setRegion(firstRegion));

            // Assuming disks are immediately available upon selecting a region
            // and regionsMap is structured to include disks directly under each region
            const disks = regionsMap[firstRegion]?.disks || [];

            // Check if there are any disks available and set related states
            if (disks.length > 0) {
                const firstDisk = disks[0];

                dispatch(setSelectedDisk(firstDisk));
            } else {
                // Handle the case where no disks are available for the first region of the selected provider
            }
        } else {
            // Handle the case where no regions (and thus no disks) are available for the selected provider
            dispatch(setRegion('')); // Clear the region or set to a default/fallback value
        }
    }

    useEffect(() => {
        nodesInRegion.forEach((node) => {
            if (node.resourceID === 0) {
                return;
            }
            nodeMap[node.resourceID] = node;
        });

        if (nodesInRegion.length > 0) {
            if (node) {
                if (!isNodeInMap(node.resourceID)) {
                    dispatch(setSelectedNode(nodesInRegion[0]));
                }
            } else {
                dispatch(setSelectedNode(nodesInRegion[0]));
            }
        }
    }, [nodeMap,node]);


    function handleAddNode(resourceID: number) {
        if (resourceID in nodeMap) {
            dispatch(setSelectedNode(nodeMap[resourceID]));
        }
    }
    function totalCost() {
        let totalBlockStorageCost = 0;
        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * disk.priceMonthly * parseInt(resource.replicas));
        }
        return node.priceMonthly * count + totalBlockStorageCost;
    }
    function totalHourlyCost() {
        let totalBlockStorageCost = 0;

        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * (disk.priceMonthly/730) * parseInt(resource.replicas));
        }
        let roundedNum = Math.ceil(node.priceHourly * Math.pow(10, 2)) / Math.pow(10, 2);
        return roundedNum * count + (totalBlockStorageCost);
    }

    return (
        <div>
            <ThemeProvider theme={mdTheme}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4}}>
                <Stack direction="row" >
                    <Card sx={{ maxWidth: 700 }}>
                        <CardContent>
                            <Typography gutterBottom variant="h5" component="div">
                                Deployment & Management
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Deploy a maximum of one app with a monthly cost up to $500/month for free. However, if a payment method is not set within one hour it will automatically delete your app.
                                You can set a payment option on the billing page. Once you've deployed an app you can view it on the clusters page within a few minutes.
                                The node sizing selection filter adds an additional 0.1 vCPU and 100Mi as overhead to help prevent selecting nodes that won't schedule this workload.
                                If a machine type you'd like isn't listed, or out of stock please contact us at support@zeus.fyi
                            </Typography>
                        </CardContent>
                        <Divider />
                        <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
                            <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between' }}>
                                <Typography variant="h6" color="text.secondary">
                                    Node Selection
                                </Typography>
                            </Box>
                            <Stack direction="column" spacing={2}>
                                <Stack direction="row" >
                                    <FormControl sx={{ mr: 2 }} fullWidth variant="outlined">
                                        <InputLabel key={`cloudProviderLabel`} id={`cloudProvider`}>
                                            Cloud Provider
                                        </InputLabel>
                                        <Select
                                            labelId={`cloudProviderLabel`}
                                            id={`cloudProvider`}
                                            name="cloudProvider"
                                            value={cloudProvider}
                                            onChange={(event) => handleChangeSelectCloudProvider(event.target.value)}
                                            label="Cloud Provider"
                                        >
                                            <MenuItem value="do">DigitalOcean</MenuItem>
                                            <MenuItem value="gcp">Google Cloud Platform</MenuItem>
                                            <MenuItem value="aws">Amazon Web Services</MenuItem>
                                            <MenuItem value="ovh">Ovh Cloud</MenuItem>
                                            <MenuItem value="azure">Azure (Coming soon)</MenuItem>
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
                                            value={region in regions ? region : ''}
                                            onChange={(event) => handleChangeSelectRegion(event.target.value)}
                                            label="Region"
                                        >
                                            {regionMenuItems}
                                        </Select>
                                    </FormControl>
                                </Stack>
                                <Stack direction="row" >
                                    {node && isNodeInMap(node.resourceID) &&
                                        <FormControl  sx={{ mr: 1 }} fullWidth variant="outlined">
                                            <InputLabel key={`nodesLabel`} id={`nodes`}>
                                                Nodes
                                            </InputLabel>
                                            <Select
                                                labelId={`nodesLabel`}
                                                id={`nodes`}
                                                name="nodes"
                                                value={node.resourceID}
                                                onChange={(event) => handleAddNode(event.target.value as number)}
                                                label="Nodes"
                                            >
                                                {nodesInRegion.map((node) => (
                                                    <MenuItem key={node.resourceID} value={node.resourceID}>
                                                        {`${node.slug} ($${node.priceMonthly.toFixed(2)}/month)`}
                                                    </MenuItem>
                                                ))}
                                            </Select>
                                        </FormControl>
                                    }
                                    <CardActions >
                                        <Stack direction="row" >
                                            <IconButton onClick={handleDecrement} aria-label="decrement" >
                                                <Remove />
                                            </IconButton>
                                            <TextField
                                                value={count}
                                                variant="outlined"
                                                size="small"
                                                inputProps={{ style: { textAlign: 'center' }, min: 0 }}
                                            />
                                            <IconButton onClick={handleIncrement} aria-label="increment">
                                                <Add />
                                            </IconButton>
                                        </Stack>
                                    </CardActions>
                                </Stack>
                                <TextField
                                    fullWidth
                                    id="description"
                                    label="Description"
                                    variant="outlined"
                                    value={node ? node.description : ""}
                                    style={{ width: "100%" }}
                                />
                                <Stack direction="row" >
                                    <TextField
                                        id="vcpus"
                                        label="vCPUs"
                                        variant="outlined"
                                        value={node ? node.vcpus : ""}
                                        sx={{ flex: 1, mr: 2 }}
                                    />
                                    <TextField
                                        id="memory"
                                        label="Memory (GB)"
                                        variant="outlined"
                                        value={node ? Math.floor(node.memory/1000) : ""}
                                        sx={{ flex: 1, mr: 2 }}
                                    />
                                    <TextField
                                        id="localDiskSize"
                                        label="Local Disk Size (GB)"
                                        variant="outlined"
                                        value={node ? node.disk : ""}
                                        sx={{ flex: 1, mr: 2 }}
                                    />
                                </Stack>
                                {node && node.gpus > 0 &&
                                    <Stack direction="row" >
                                        <TextField
                                            id="gpuType"
                                            label="gpuType"
                                            variant="outlined"
                                            value={node.gpuType}
                                            sx={{ flex: 1, mr: 2 }}
                                        />
                                        <TextField
                                            id="gpus"
                                            label="gpus"
                                            variant="outlined"
                                            value={node.gpus}
                                            sx={{ flex: 1, mr: 2 }}
                                        />
                                    </Stack>
                                }
                                <Divider />
                                <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="h6" color="text.secondary">
                                        Block Storage
                                    </Typography>
                                </Box>
                                {resourceRequirements.map((resourceRequirement, index) => (
                                    <div key={index}>
                                        <Stack direction="row" >
                                            <TextField
                                                fullWidth
                                                id={`componentName-${index}`}
                                                label="Cluster Base"
                                                variant="outlined"
                                                value={resourceRequirement.componentBaseName}
                                                sx={{ flex: 1, mr: 2 }}
                                            />
                                            <TextField
                                                fullWidth
                                                id={`blockStorageSize-${index}`}
                                                label="PVC Disk Size SSD"
                                                variant="outlined"
                                                value={resourceRequirement.resourceSumsDisk}
                                                sx={{ flex: 1, mr: 2 }}
                                            />
                                            <TextField
                                                value={resourceRequirement.replicas}
                                                fullWidth
                                                id={`replicas-${index}`}
                                                label="Replicas"
                                                variant="outlined"
                                                sx={{ flex: 1, mr: 2 }}
                                            />
                                        </Stack>
                                    </div>
                                ))}
                                {
                                    (
                                        disksInRegion.length > 0 &&
                                        <div>
                                            <Typography variant="h6" color="text.secondary">
                                                Disk Type and Pricing
                                            </Typography>
                                            <Box sx={{mt: 2, mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                <Select
                                                    labelId="disk-type-select-label"
                                                    id="disk-type-select"
                                                    value={disk && disk.type + "-" + disk.subType || disksInRegion[0].type+'-'+disksInRegion[0].subType}
                                                    label="Disk Type"
                                                    sx={{ width: 200, mr: 2 }}
                                                    onChange={handleDiskTypeChange}
                                                >
                                                    {diskCloudRegionMenuItems}
                                                </Select>
                                                <TextField
                                                    value={disk.priceMonthly.toFixed(2)}
                                                    fullWidth
                                                    id={`monthlyPrice-${disk.priceMonthly}`}
                                                    label="Monthly Cost ($)"
                                                    variant="outlined"
                                                    inputProps={{ readOnly: true }}
                                                    sx={{ flex: 1, mr: 2 }}
                                                />
                                                <TextField
                                                    value={(disk.priceMonthly/730).toFixed(4)}
                                                    fullWidth
                                                    id={`hourlyPrice-${disk.priceHourly}`}
                                                    label="Hourly Cost ($)"
                                                    variant="outlined"
                                                    inputProps={{ readOnly: true }}
                                                    sx={{ flex: 1, mr: 2 }}
                                                />
                                            </Box>
                                        </div>

                                    )
                                }
                                <Divider />
                                <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between' }}>
                                    <Typography variant="h6" color="text.secondary">
                                        Total Costs
                                    </Typography>
                                </Box>
                                <Stack direction="row" >
                                    <TextField
                                        fullWidth
                                        id="monthlyCost"
                                        label="Monthly Cost ($)"
                                        variant="outlined"
                                        value={node ? totalCost().toFixed(2) : ""}
                                        sx={{ flex: 1, mr: 2 }}
                                    />
                                    <TextField
                                        fullWidth
                                        id="hourlyCost"
                                        label="Hourly Cost ($)"
                                        variant="outlined"
                                        value={node ? totalHourlyCost().toFixed(4) : ""}
                                        sx={{ flex: 1, mr: 2 }}
                                    />
                                    <CardActions >
                                        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                            <Button variant="contained" onClick={handleDeploy} disabled={buttonDisabled}>{buttonLabel}</Button>
                                        </div>
                                    </CardActions>
                                </Stack>
                                <Box >
                                    {statusMessage && (
                                        <Typography variant="body2" color={requestStatus === 'error' || requestStatus === 'missingBilling' || requestStatus === 'outOfCredits' ? 'error' : 'success'}>
                                            {statusMessage}
                                        </Typography>
                                    )}
                                </Box>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Requested Resource Summary
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Reads your Kubernetes workloads and summarizes the requested resources to deploy this app.
                                    </Typography>
                                </CardContent>
                                <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
                                    <Box sx={{ mt: 2, display: 'flex' }}>
                                        <ResourceRequirementsTable />
                                    </Box>
                                </Container>
                            </Stack>
                        </Container>
                    </Card>
                    {cluster.clusterName && cluster.clusterName === 'sui' &&
                        <Card sx={{ width: '20%' }}>
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="div">
                                    Config Options
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Select which components and config options you want to deploy for this app. This will create a copy of the config into your private apps
                                    and will automatically navigate to the app deployment page for that configuration.
                                </Typography>
                            </CardContent>
                            <AppConfigsTable />
                        </Card>
                    }
                </Stack>
            </Container>
            </ThemeProvider>
        </div>
    );
}


function getSuffix(input: string): string {
    if (input.length === 0) {
        return '';
    }
    const parts = input.split('-');
    return parts[parts.length - 1];
}

function getPrefix(input: string): string {
    if (input.length === 0) {
        return '';
    }
    const parts = input.split('-');
    return parts[0];
}
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
    setCluster,
    setClusterPreview,
    setNodes,
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../redux/apps/apps.reducer";

const mdTheme = createTheme();

interface NodeMap {
    [resourceID: number]: Nodes;
}

export function DeployPage(props: any) {
    const {app, region, setRegion, cloudProvider, setCloudProvider} = props
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const resourceRequirements = createDiskResourceRequirements(cluster);
    let nodes = useSelector((state: RootState) => state.apps.nodes);
    const nodeMap: NodeMap = {};
    const [count, setCount] = useState(0);
    const [freeTrial, setFreeTrial] = useState(false);
    let filteredNodes = nodes.filter((node) => node.cloudProvider === cloudProvider && node.region === region);
    const [node, setNode] = useState(filteredNodes[0]);
    const params = useParams();
    const dispatch = useDispatch();

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
                if (response.nodes.length > 0) {
                    dispatch(setNodes(response.nodes))
                }
                nodes = response.nodes
                filteredNodes = nodes.filter((node) => node.cloudProvider === cloudProvider && node.region === region);
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

        if (filteredNodes.length > 0 && filteredNodes[0].resourceID === 0) {
            fetchData().then(r => {
                setNode(filteredNodes[0]);
            });
        }
    }, [params.id, nodes, filteredNodes, nodeMap, cloudProvider, region]);

    const handleIncrement = () => {
        setCount(count + 1);
    };

    const handleDecrement = () => {
        if (count - 1 < 0) {
            setCount(0)
            return;
        }
        setCount(count - 1);
    };

    filteredNodes.forEach((node) => {
        if (node.resourceID === 0) {
            return;
        }
        nodeMap[node.resourceID] = node;
    });

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
            buttonLabel = 'Deploy';
            buttonDisabled = false;
            break;
    }
    const handleDeploy = async () => {
        try {
            setRequestStatus('pending');
            const namespaceAlias = cluster.clusterName;
            const payload = {
                "cloudProvider": cloudProvider,
                "region": region,
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
    function handleChangeSelectCloudProvider(cloudProvider: string) {
        setCloudProvider(cloudProvider);
        if (cloudProvider === 'gcp') {
            setRegion('us-central1');
        }
        if (cloudProvider === 'do') {
            setRegion('nyc1');
        }
        if (cloudProvider === 'aws') {
            setRegion('us-west-1');
        }
        if (cloudProvider === 'ovh') {
            setRegion('us-west-or-1');
        }
    }
    useEffect(() => {
        filteredNodes = nodes.filter((node) => node.cloudProvider === cloudProvider && node.region === region);
        filteredNodes.forEach((node) => {
            if (node.resourceID === 0) {
                return;
            }
            nodeMap[node.resourceID] = node;
        });

        if (filteredNodes.length > 0) {
            if (node) {
                if (!isNodeInMap(node.resourceID)) {
                    setNode(filteredNodes[0]);
                }
            } else {
                setNode(filteredNodes[0]);
            }
        }
    }, [cloudProvider, region, nodeMap,node]);

    function handleChangeSelectRegion(region: string) {
        if (cloudProvider === 'aws') {
            setRegion('us-west-1');
        } else if (cloudProvider === 'gcp') {
            setRegion('us-central1')
        } else if (cloudProvider == 'ovh') {
            setRegion('us-west-or-1');
        } else {
            setRegion(region);
        }
    }

    function isNodeInMap(resourceID: number) {
        return resourceID in nodeMap;
    }
    function handleAddNode(resourceID: number) {
        if (resourceID in nodeMap) {
            setNode(nodeMap[resourceID]);
        }
    }
    function totalCost() {
        let totalBlockStorageCost = 0;
        // digitalOcean block storage
        let monthlyDiskCost = 10
        if (cloudProvider === 'gcp') {
            monthlyDiskCost = 17
        }
        if (cloudProvider === 'aws') {
            monthlyDiskCost = 12.88
        }
        if (cloudProvider === 'ovh') {
            monthlyDiskCost = 12
        }
        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * monthlyDiskCost * parseInt(resource.replicas));
        }
        return node.priceMonthly * count + (totalBlockStorageCost*1.1);
    }
    function totalHourlyCost() {
        let totalBlockStorageCost = 0;
        // digitalOcean block storage
        let hourlyDiskCost = 0.0137
        if (cloudProvider === 'gcp') {
            hourlyDiskCost = 0.02329
        }
        if (cloudProvider === 'aws') {
            hourlyDiskCost = 0.01765
        }
        if (cloudProvider == 'ovh') {
            hourlyDiskCost = 0.01643835616
        }
        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * hourlyDiskCost * parseInt(resource.replicas));
        }
        let roundedNum = Math.ceil(node.priceHourly * Math.pow(10, 2)) / Math.pow(10, 2);
        return roundedNum * count + (totalBlockStorageCost*1.1);
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
                            Without setting up a payment method you can only deploy a maximum of one app with a monthly cost up to $500/month, and if a payment method is not set within one hour it will automatically delete your app.
                            You can set a payment option on the billing page. Once you've deployed an app you can view it on the clusters page within a few minutes. Click on the cluster namespace to get a detailed view of the live cluster.
                            The node sizing selection filter adds an additional 0.1 vCPU and 1.5Gi as overhead from the server to prevent selecting nodes that won't schedule this workload.
                            If a machine type you'd like isn't listed please contact us at alex@zeus.fyi
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
                                        value={region}
                                        onChange={(event) => handleChangeSelectRegion(event.target.value)}
                                        label="Region"
                                    >
                                        {
                                            (() => {
                                                switch (cloudProvider) {
                                                    case 'gcp':
                                                        return <MenuItem value="us-central1">us-central1</MenuItem>;
                                                    case 'aws':
                                                        return <MenuItem value="us-west-1">us-west-1</MenuItem>; // Add the respective region for AWS
                                                    case 'ovh':
                                                        return <MenuItem value="us-west-or-1">us-west-or-1</MenuItem>;
                                                    default:
                                                        return <MenuItem value="nyc1">nyc1</MenuItem>; // Default is for any other provider
                                                }
                                            })()
                                        }
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
                                        {nodes
                                            .filter((node) => node.cloudProvider === cloudProvider && node.region === region)
                                            .map((node) => (
                                            <MenuItem key={node.resourceID} value={node.resourceID}>
                                                {node.slug + ' ($' + node.priceMonthly.toFixed(2) + '/month)'}
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
                                    value={node ? totalHourlyCost().toFixed(2) : ""}
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
            </ThemeProvider>

        </div>
    );
}
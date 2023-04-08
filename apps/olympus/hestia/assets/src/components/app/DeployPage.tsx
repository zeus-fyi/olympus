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
    [nodeId: number]: Nodes;
}

export function DeployPage(props: any) {
    const [cloudProvider, setCloudProvider] = useState('do');
    const [region, setRegion] = useState('nyc1');
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const resourceRequirements = createDiskResourceRequirements(cluster);
    let nodes = useSelector((state: RootState) => state.apps.nodes);
    const nodeMap: NodeMap = {};
    const [count, setCount] = useState(0);
    const [node, setNode] = useState(nodes[0]);
    const params = useParams();
    const dispatch = useDispatch();

    useEffect(() => {
        async function fetchData() {
            let cluster;
            let clusterPreview;
            let selectedComponentBaseName;
            let selectedSkeletonBaseName;
            try {
                const response = await appsApiGateway.getPrivateAppDetails(params.id as string);
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
                return response;
            } catch (e) {
            }
        }
        if (nodes[0].nodeID === 0) {
            fetchData().then(r => {
                setNode(nodes[0]);
            });
        }

    }, [params.id, nodes]);

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

    nodes.forEach((node) => {
        if (node.nodeID === 0) {
            return;
        }
        nodeMap[node.nodeID] = node;
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

    function handleChangeSelectCloudProvider(cloudProvider: string) {
        setCloudProvider(cloudProvider);
    }

    function handleChangeSelectRegion(region: string) {
        setRegion(region);
    }

    function isNodeInMap(nodeID: number) {
        return nodeID in nodeMap;
    }

    function handleAddNode(nodeID: number) {
        if (nodeID in nodeMap) {
            setNode(nodeMap[nodeID]);
        }
    }
    function totalCost() {
        let totalBlockStorageCost = 0;
        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * 10 * parseInt(resource.replicas));
        }
        return node.priceMonthly * count + totalBlockStorageCost;
    }
    function totalHourlyCost() {
        let totalBlockStorageCost = 0;
        for (const resource of resourceRequirements) {
            totalBlockStorageCost += (Number(resource.blockStorageCostUnit) * 0.10 * parseInt(resource.replicas));
        }
        return node.priceHourly * count + totalBlockStorageCost;
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
                            Currently, you cannot deploy clusters without getting authorization manually first until we have automated billing setup.
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
                            <Stack direction="row" >
                                {isNodeInMap(node.nodeID) &&
                                <FormControl  sx={{ mr: 1 }} fullWidth variant="outlined">
                                    <InputLabel key={`nodesLabel`} id={`nodes`}>
                                        Nodes
                                    </InputLabel>
                                    <Select
                                        labelId={`nodesLabel`}
                                        id={`nodes`}
                                        name="nodes"
                                        value={node.nodeID}
                                        onChange={(event) => handleAddNode(event.target.value as number)}
                                        label="Nodes"
                                    >
                                        {nodes.map((node) => (
                                            <MenuItem key={node.nodeID} value={node.nodeID}>
                                                {node.slug}
                                            </MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                                }
                                <TextField
                                    fullWidth
                                    id="description"
                                    label="Description"
                                    variant="outlined"
                                    value={node ? node.description : ""}
                                    style={{ width: "50%" }}
                                />
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
                            </Stack>
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
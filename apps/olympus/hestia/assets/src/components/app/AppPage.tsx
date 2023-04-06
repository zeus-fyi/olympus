import * as React from "react";
import {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import {AppBar, Drawer} from "../dashboard/Dashboard";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import MenuIcon from "@mui/icons-material/Menu";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import Divider from "@mui/material/Divider";
import List from "@mui/material/List";
import MainListItems from "../dashboard/listItems";
import Container from "@mui/material/Container";
import {appsApiGateway} from "../../gateway/apps";
import {RootState} from "../../redux/store";
import {setSelectedComponentBaseName, setSelectedSkeletonBaseName} from "../../redux/apps/apps.reducer";
import TextField from "@mui/material/TextField";
import {Card, CardContent, FormControl, InputLabel, MenuItem, Select, Stack} from "@mui/material";
import {setClusterPreview, setSelectedContainerName} from "../../redux/clusters/clusters.builder.reducer";
import YamlTextFieldAppPage from "./YamlFormattedTextAppPage";

const mdTheme = createTheme();

function createTopologyData(
    topologyID: number,
    clusterName: string,
    componentBaseName: string,
    skeletonBaseName: string,
) {
    return {topologyID, clusterName, componentBaseName, skeletonBaseName};
}

function AppPageContent() {
    const [open, setOpen] = React.useState(true);
    const [previewType, setPreviewType] = useState('');
    const [addDeployment, setAddDeployment] = useState(true);
    const [addConfigMap, setAddConfigMap] = useState(true);
    const [addIngress, setAddIngress] = useState(true);
    const [addService, setAddService] = useState(true);
    const [addStatefulSet, setAddStatefulSet] = useState(true);

    const params = useParams();
    const dispatch = useDispatch();
    let clusterPreview = useSelector((state: RootState) => state.apps.clusterPreview);
    const [name, setName] = useState('');
    useEffect(() => {
        async function fetchData() {
            try {
                const response = await appsApiGateway.getPrivateAppDetails(params.id as string);
                dispatch(setClusterPreview(response.clusterPreview));
                setName(response.clusterPreview.clusterName);
                const cb = Object.keys(response.cluster.componentBases)
                if (cb.length > 0) {
                    dispatch(setSelectedComponentBaseName(cb[0]));
                    const sbs = Object.keys(response.cluster.componentBases[cb[0]])
                    if (sbs.length > 0) {
                        dispatch(setSelectedSkeletonBaseName(sbs[0]));
                    }
                }
                return response;
            } catch (e) {
            }
        }
        fetchData();
    }, [params.id, dispatch]);
    const onClickView = (newPreviewType: string) => {
        setPreviewType(newPreviewType);
    }
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }
    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open} style={{ backgroundColor: '#8991B0'}}>
                    <Toolbar
                        sx={{
                            pr: '24px', // keep right padding when drawer closed
                        }}
                    >
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && { display: 'none' }),
                            }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{ flexGrow: 1 }}
                        >
                            Clusters
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={handleLogout}
                        >Logout
                        </Button>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeftIcon />
                        </IconButton>
                    </Toolbar>
                    <Divider />
                    <List component="nav">
                        <MainListItems />
                        <Divider sx={{ my: 1 }} />
                    </List>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: (theme) =>
                            theme.palette.mode === 'light'
                                ? theme.palette.grey[100]
                                : theme.palette.grey[900],
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                    }}
                >
                    <Toolbar />
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Stack spacing={2} direction="row">
                                <Card sx={{ maxWidth: 500 }}>
                                    <CardContent>
                                        <Typography gutterBottom variant="h5" component="div">
                                            Cluster App Details
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary">
                                            Review Cluster Class & Workload Bases
                                        </Typography>
                                    </CardContent>
                                    <Container maxWidth="xl" sx={{ mb: 4 }}>
                                        <Box mt={2}>
                                            <TextField
                                                fullWidth
                                                id={`clusterName`}
                                                label={`Cluster Name`}
                                                variant="outlined"
                                                value={name}
                                                InputProps={{ readOnly: true }}
                                            />
                                        </Box>
                                        <Box mt={2}>
                                            <SelectedComponentBaseName />
                                        </Box>
                                        <Box mt={2}>
                                            <SelectedSkeletonBaseName />
                                        </Box>
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
                                <YamlTextFieldAppPage previewType={previewType} clusterPreview={clusterPreview}/>
                            </Stack>
                        </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}


export function SelectedComponentBaseName(props: any) {
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.apps.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.apps.selectedComponentBaseName);
    const onAccessComponentBase = (selectedComponentBaseName: string) => {
        dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
        const skeletonBaseName = Object.keys(cluster.componentBases[selectedComponentBaseName])[0];
        dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
        // Add a check to see if the `containers` field exists
        if (cluster.componentBases[selectedComponentBaseName] &&
            cluster.componentBases[selectedComponentBaseName][skeletonBaseName] &&
            cluster.componentBases[selectedComponentBaseName][skeletonBaseName].containers) {
            const containerKeys = Object.keys(cluster.componentBases[selectedComponentBaseName][skeletonBaseName].containers);
            if (containerKeys.length > 0) {
                dispatch(setSelectedContainerName(containerKeys[0]));
            }
        }
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

export function SelectedSkeletonBaseName(props: any) {
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
export default function AppPage() {
    return <AppPageContent />
}
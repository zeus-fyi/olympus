import * as React from 'react';
import {useEffect, useState} from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import {AppBar, Drawer} from '../dashboard/Dashboard';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import {Card, CardContent, createTheme, TableContainer, TableRow} from '@mui/material';
import TableBody from "@mui/material/TableBody";
import MainListItems from "../dashboard/listItems";
import {clustersApiGateway} from "../../gateway/clusters";
import {ThemeProvider} from "@mui/material/styles";
import Stack from "@mui/material/Stack";
import {CloudCtxNs, resourcesApiGateway} from "../../gateway/resources";
import {ClusterViews} from "./ClusterAppViews";
import {RootState} from "../../redux/store";
import {
    setClusterViewEnabledToggle,
    setSelectedClusterAppViewName
} from "../../redux/clusters/clusters.builder.reducer";

const mdTheme = createTheme();

function createData(
    cloudCtxNsID: number,
    cloudProvider: string,
    region: string,
    context: string,
    namespace: string,
    namespaceAlias: string,
) {
    return {cloudCtxNsID, cloudProvider, region, context, namespace, namespaceAlias};
}

function createClusterAppViewData(
    cloudCtxNsID: number,
    clusterClassName: string,
    cloudProvider: string,
    region: string,
    context: string,
    namespace: string,
    namespaceAlias: string,
    appName: string,
) {
    if (appName !== '' && clusterClassName == appName) {
        return {cloudCtxNsID, clusterClassName, cloudProvider, region, context, namespace, namespaceAlias};
    }
    if (appName === '') {
        return {cloudCtxNsID, clusterClassName, cloudProvider, region, context, namespace, namespaceAlias};
    }
    return {}
}
function ClustersContent() {
    const [open, setOpen] = React.useState(true);
    const appName = useSelector((state: RootState) => state.clusterBuilder.selectedClusterAppView);
    const pageView = useSelector((state: RootState) => state.clusterBuilder.clusterViewEnabledToggle);
    const [clusters, setClusters] = useState([{}]);
    const [allClusters, setAllClusters] = useState([{}]);
    const dispatch = useDispatch();
    let navigate = useNavigate();

    const setPageView = (pageView: boolean) => {
        dispatch(setClusterViewEnabledToggle(pageView));
    }

    const toggleDrawer = () => {
        setOpen(!open);
    }

    const setAppName = (appName: string) => {
        dispatch(setSelectedClusterAppViewName(appName));
    }

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
                <AppBar position="absolute" open={open} style={{ backgroundColor: '#333'}}>
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
                        <Stack direction="row" alignItems="center" spacing={2}>
                            <Card sx={{ maxWidth: 600 }}>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Cloud, Infra, & Cluster Management
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                       Clusters created will show up here after shortly after app deployment. Click the cluster row to get a detailed view of the cluster.
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Stack>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Card sx={{ maxWidth: 600 }}>
                            <ClusterViews pageView={pageView} setPageView={setPageView} appName={appName} setAppName={setAppName} clusters={clusters} allClusters={allClusters}/>
                        </Card>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        {pageView ? <CloudClustersAppsView clusters={clusters} setClusters={setClusters} appName={appName} setAllClusters={setAllClusters} /> : <CloudClusters />}
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

function ClustersTable(clusters: any) {
    let navigate = useNavigate();
    const [cluster, setCluster] = useState([{}]);
    const [statusMessage, setStatusMessage] = useState('');
    const [statusMessageRowIndex, setStatusMessageRowIndex] = useState(-1);

    const handleClick = async (event: any, cluster: any) => {
        const tableRow = event.currentTarget;
        const tableCells = tableRow.children;
        if (event.target === tableCells[tableCells.length - 1] || event.target === tableCells[4] || event.target === tableCells[5]) {
            return;
        }
        event.preventDefault();
        setCluster(cluster);
        navigate('/clusters/'+cluster.cloudCtxNsID);
    }

    const handleDeleteNamespace = async (index: number, cloudCtxNsId: number, cloudProvider: string, region: string, context: string, namespace: string) => {
        try {
            const cloudCtxNs = {
                cloudProvider: cloudProvider,
                region: region,
                context: context,
                namespace: namespace,
            } as CloudCtxNs;
            const response = await resourcesApiGateway.destroyDeploy(cloudCtxNs);
            setStatusMessage(`Destroy in progress`);
            setStatusMessageRowIndex(index);
        } catch (error) {
            console.error(error);
            setStatusMessage(`Error deleting cloudCtxNs ID ${cloudCtxNsId}`);
            setStatusMessageRowIndex(index);
        }
    }
    return( <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
                <TableRow style={{ backgroundColor: '#333'}} >
                    <TableCell style={{ color: 'white'}}>CloudCtxNsID</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">CloudProvider</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Region</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Context</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Namespace</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Alias</TableCell>
                    <TableCell style={{ color: 'white'}} align="left"></TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {clusters.map((row: any, i: number) => (
                    <TableRow
                        key={i}
                        onClick={(event) => handleClick(event, row)}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                        <TableCell component="th" scope="row">
                            {row.cloudCtxNsID}
                        </TableCell>
                        <TableCell align="left">{row.cloudProvider}</TableCell>
                        <TableCell align="left">{row.region}</TableCell>
                        <TableCell align="left">{row.context}</TableCell>
                        <TableCell align="left">{row.namespace}</TableCell>
                        <TableCell align="left">{row.namespaceAlias}</TableCell>
                        <TableCell align="left">
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={async (event) => {
                                    event.stopPropagation();
                                    await handleDeleteNamespace(i,row.cloudCtxNsID, row.cloudProvider, row.region, row.context, row.namespace);
                                }}
                            >
                                Delete
                            </Button>
                            {statusMessageRowIndex === i && <div>{statusMessage}</div>}
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    </TableContainer>)
}

function CloudClusters() {
    const [clusters, setClusters] = useState([{}]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await clustersApiGateway.getClusters();
                const clustersData: any[] = response.data;
                const clusterRows = clustersData.map((cluster: any) =>
                    createData(cluster.cloudCtxNsID, cluster.cloudProvider, cluster.region, cluster.context, cluster.namespace, cluster.namespaceAlias),
                );
                setClusters(clusterRows)
            } catch (error) {
                console.log("error", error);
            }}
        fetchData().then(r => '');
    }, []);
    return (
        ClustersTable(clusters)
    )
}

export default function Clusters(props: any) {
    return <ClustersContent />
}

function CloudClustersAppsView(props: any) {
    const { appName, clusters, setClusters, setAllClusters } = props;

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await clustersApiGateway.getAppClustersView();
                const clustersData: any[] = response.data;
                const allClusters = clustersData.map((cluster: any) =>
                    createClusterAppViewData(cluster.cloudCtxNsID, cluster.clusterClassName,cluster.cloudProvider, cluster.region, cluster.context, cluster.namespace, cluster.namespaceAlias, ''),
                )
                setAllClusters(allClusters);
                const clusterRows = clustersData.map((cluster: any) =>
                    createClusterAppViewData(cluster.cloudCtxNsID, cluster.clusterClassName,cluster.cloudProvider, cluster.region, cluster.context, cluster.namespace, cluster.namespaceAlias, appName),
                ).filter((clusterViewData: any) => Object.keys(clusterViewData).length !== 0);
                setClusters(clusterRows)
            } catch (error) {
                console.log("error", error);
            }}
        fetchData().then(r => '');
    }, [appName]);
    return (
        CloudClustersAppViewTable(clusters)
    )
}

function CloudClustersAppViewTable(clusters: any) {
    let navigate = useNavigate();
    const [cluster, setCluster] = useState([{}]);
    const [statusMessage, setStatusMessage] = useState('');
    const [statusMessageRowIndex, setStatusMessageRowIndex] = useState(-1);

    const handleClick = async (event: any, cluster: any) => {
        const tableRow = event.currentTarget;
        const tableCells = tableRow.children;
        if (event.target === tableCells[tableCells.length - 1]) {
            return;
        }
        event.preventDefault();
        setCluster(cluster);
        navigate('/clusters/'+cluster.cloudCtxNsID);
    }

    const handleDeleteNamespace = async (index: number, cloudCtxNsId: number, cloudProvider: string, region: string, context: string, namespace: string) => {
        try {
            const cloudCtxNs = {
                cloudProvider: cloudProvider,
                region: region,
                context: context,
                namespace: namespace,
            } as CloudCtxNs;
            const response = await resourcesApiGateway.destroyDeploy(cloudCtxNs);
            setStatusMessage(`Destroy in progress`);
            setStatusMessageRowIndex(index);
        } catch (error) {
            console.error(error);
            setStatusMessage(`Error deleting cloudCtxNs ID ${cloudCtxNsId}`);
            setStatusMessageRowIndex(index);
        }
    }
    return( <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
                <TableRow style={{ backgroundColor: '#333'}} >
                    <TableCell style={{ color: 'white'}}>CloudCtxNsID</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">AppName</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">CloudProvider</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Region</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Context</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Namespace</TableCell>
                    <TableCell style={{ color: 'white'}} align="left">Alias</TableCell>
                    <TableCell style={{ color: 'white'}} align="left"></TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {clusters.map((row: any, i: number) => (
                    <TableRow
                        key={i}
                        onClick={(event) => handleClick(event, row)}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                        <TableCell component="th" scope="row">
                            {row.cloudCtxNsID}
                        </TableCell>
                        <TableCell align="left">{row.clusterClassName}</TableCell>
                        <TableCell align="left">{row.cloudProvider}</TableCell>
                        <TableCell align="left">{row.region}</TableCell>
                        <TableCell align="left">{row.context}</TableCell>
                        <TableCell align="left">{row.namespace}</TableCell>
                        <TableCell align="left">{row.namespaceAlias}</TableCell>
                        <TableCell align="left">
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={async (event) => {
                                    event.stopPropagation();
                                    await handleDeleteNamespace(i,row.cloudCtxNsID, row.cloudProvider, row.region, row.context, row.namespace);
                                }}
                            >
                                Delete
                            </Button>
                            {statusMessageRowIndex === i && <div>{statusMessage}</div>}
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    </TableContainer>)
}

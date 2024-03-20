import * as React from "react";
import {useEffect, useState} from "react";
import {clustersApiGateway} from "../../gateway/clusters";
import {Card, CardContent, TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import {useNavigate, useParams} from "react-router-dom";
import {useDispatch} from "react-redux";
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
import {PodsPageTable} from "./Pods";

const mdTheme = createTheme();

function createTopologyData(
    topologyID: number,
    clusterName: string,
    componentBaseName: string,
    skeletonBaseName: string,
) {
    return {topologyID, clusterName, componentBaseName, skeletonBaseName};
}

function ClustersPageContent() {
    const [open, setOpen] = React.useState(true);

    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();

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
                        <Card>
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="div">
                                    Cluster Details
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    You can see the details of the cluster here, and interact, debug, and rollout changes.
                                </Typography>
                            </CardContent>
                        </Card>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <ClustersPageTable />
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <PodsPageTable />
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

function ClustersPageTable(cluster: any) {
    const params = useParams();
    const [activeClusterTopologies, setActiveClusterTopologies] = useState([{}]);
    const [statusMessage, setStatusMessage] = useState('');
    const [statusMessageRowIndex, setStatusMessageRowIndex] = useState<number | null>(null);
    let navigate = useNavigate();

    //                 navigate('/signup');

    const onClickAppsPage = async (index: number, appName: any) => {
        //event.preventDefault();
        navigate('/app/'+appName);
        console.log("appName", appName)
    }
    const onClickRolloutUpgrade = async (index: number, clusterClassName: string) => {
        try {
            const response = await clustersApiGateway.deployClusterToCloudCtxNs(params.id, clusterClassName, activeClusterTopologies);
            const statusCode = response.status;
            if (statusCode === 202) {
                setStatusMessageRowIndex(index);
                setStatusMessage(`Cluster ${clusterClassName} update in progress`);
            } else if (statusCode === 200){
                setStatusMessageRowIndex(index);
                setStatusMessage(`Cluster ${clusterClassName} already up to date`);
            } else {
                setStatusMessageRowIndex(index);
                setStatusMessage(`Cluster ${clusterClassName} had an unexpected response: status code ${statusCode}`);
            }
        } catch (e) {
            setStatusMessageRowIndex(index);
            setStatusMessage(`Cluster ${clusterClassName} failed to update`);
        }
    }

    const onClickRolloutRestart = async (index: number, clusterClassName: string) => {
        try {
            const response = await clustersApiGateway.deployRolloutRestartApp(params.id);
            const statusCode = response.status;
            if (statusCode < 400) {
                setStatusMessageRowIndex(index);
                setStatusMessage(`Cluster ${clusterClassName} restart in progress`);
            } else {
                setStatusMessageRowIndex(index);
                setStatusMessage(`Cluster ${clusterClassName} had an unexpected response: status code ${statusCode}`);
            }
        } catch (e) {
            setStatusMessageRowIndex(index);
            setStatusMessage(`Cluster ${clusterClassName} failed to update`);
        }
    }

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await clustersApiGateway.getClusterTopologies(params);
                const clustersTopologyData: any[] = response.data;
                const clusterTopologyRows = clustersTopologyData.map((topology: any) =>
                    createTopologyData(topology.topologyID, topology.clusterName, topology.componentBaseName, topology.skeletonBaseName),
                );
                setActiveClusterTopologies(clusterTopologyRows);
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(params);
    }, []);
    return (
        <div>
        <Box sx={{ mt: 4, mb: 4 }}>
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                    <TableHead>
                        <TableRow style={{ backgroundColor: '#333'}} >
                            <TableCell style={{ color: 'white'}} align="left">ClusterName</TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            <TableCell style={{ color: 'white'}} align="left"></TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {activeClusterTopologies
                            .filter(
                                (item: any, index: number, self: any) =>
                                    index ===
                                    self.findIndex(
                                        (otherItem: any) => otherItem.clusterName === item.clusterName
                                    )
                            )
                            .map((row: any, i: number) => (
                            <TableRow
                                key={i}
                                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                            >
                                <TableCell component="th" scope="row">
                                    {row.clusterName}
                                </TableCell>
                                <TableCell align="left">
                                    <Button onClick={() => onClickAppsPage(i, row.clusterName)} variant="contained">zK8s App Page</Button>
                                </TableCell>
                                <TableCell align="left">
                                    <Button onClick={() => onClickRolloutUpgrade(i, row.clusterName)} variant="contained">Deploy Latest</Button>
                                </TableCell>
                                <TableCell align="left">
                                    <Button onClick={() => onClickRolloutRestart(i, row.clusterName)} variant="contained">Rollout Restart</Button>
                                </TableCell>
                                <TableCell align="left">
                                    The Deploy Latest button will deploy the latest version configs to the cluster. The Rollout Restart button will restart the cluster.
                                </TableCell>
                                <Box sx={{ml: 2, mt: 2}}>
                                    {statusMessageRowIndex === i && <div>{statusMessage}</div>}
                                </Box>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ color: 'white'}}>TopologyID</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">ClusterName</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">ClusterBaseName</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">SkeletonBaseName</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {activeClusterTopologies.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.topologyID}
                            </TableCell>
                            <TableCell align="left">{row.clusterName}</TableCell>
                            <TableCell align="left">{row.componentBaseName}</TableCell>
                            <TableCell align="left">{row.skeletonBaseName}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
        </div>
    );
}

export default function ClustersPage() {
    return <ClustersPageContent />
}
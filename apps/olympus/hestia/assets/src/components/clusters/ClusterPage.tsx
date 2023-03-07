import * as React from "react";
import {useEffect, useState} from "react";
import {clustersApiGateway} from "../../gateway/clusters";
import {TableContainer, TableRow} from "@mui/material";
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

    const handleLogout = (event: any) => {
        event.preventDefault();
        authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }
    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open}>
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
                            Dashboard
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
                    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                        <ClustersPageTable />
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

function ClustersPageTable(cluster: any) {
    const params = useParams();
    const [activeClusterTopologies, setActiveClusterTopologies] = useState([{}]);
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await clustersApiGateway.getClusterTopologies(parseInt(params.id));
                console.log(response.data)
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
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>TopologyID</TableCell>
                        <TableCell align="left">ClusterName</TableCell>
                        <TableCell align="left">ClusterBaseName</TableCell>
                        <TableCell align="left">SkeletonBaseName</TableCell>
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
    );
}

export default function ClustersPage() {
    return <ClustersPageContent />
}
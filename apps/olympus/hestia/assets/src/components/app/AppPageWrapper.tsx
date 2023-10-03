import {createTheme, ThemeProvider} from "@mui/material/styles";
import * as React from "react";
import {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
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
import {
    setCluster,
    setClusterPreview,
    setNodes,
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../redux/apps/apps.reducer";
import DeployConfigToggle from "./DeployConfigToggle";

const mdTheme = createTheme();

export function AppPageWrapper(props: any) {
    const {app} = props
    const [cloudProvider, setCloudProvider] = useState('do');
    const [region, setRegion] = useState('nyc1');
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let clusterPreview = useSelector((state: RootState) => state.apps.clusterPreview);
    let cluster = useSelector((state: RootState) => state.apps.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.apps.selectedComponentBaseName);
    let selectedSkeletonBaseName = useSelector((state: RootState) => state.apps.selectedSkeletonBaseName);
    let navigate = useNavigate();
    const dispatch = useDispatch();
    const params = useParams();
    const handleLogout = (event: any) => {
        event.preventDefault();
        authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    useEffect(() => {
        async function fetchData() {
            try {
                let id = params.id as string;
                if (app === "avax") {
                    id = "avax"
                }
                if (app === "ethereumEphemeralBeacons" ){
                    id = 'ethereumEphemeralBeacons'
                }
                if (app === "microservice" ){
                    id = 'microservice'
                }
                if (app === "sui" ){
                    id = 'sui'
                }
                console.log('gdssdfsdf')
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
                return response;
            } catch (e) {
                //console.log(e, 'error')
            }
        }
        fetchData().then(r => {
        });
    }, [params.id, app]);
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
                            Infra Builder
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
                    <div style={{ display: 'flex' }}>
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <DeployConfigToggle app={app} cloudProvider={cloudProvider} setCloudProvider={setCloudProvider} region={region} setRegion={setRegion}/>
                        </Container>
                    </div>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

import * as React from "react";
import {useEffect} from "react";
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
import {setSelectedClusterApp} from "../../redux/apps/apps.reducer";

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
                    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                        <AppPageDetails />
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

function AppPageDetails(props: any) {
    const params = useParams();
    const selectedApp = useSelector((state: RootState) => state.apps.selectedClusterApp);
    const dispatch = useDispatch();

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await appsApiGateway.getPrivateAppDetails(params.id as string);
                dispatch(setSelectedClusterApp(response));
            } catch (e) {
            }
        }
        fetchData();
    }, [dispatch]);

    const name = selectedApp.clusterName
    return (
        <div>{name}</div>
    );
}

export default function AppPage() {
    return <AppPageContent />
}
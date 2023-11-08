import * as React from 'react';
import {createTheme, styled, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, {AppBarProps as MuiAppBarProps} from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {Card, CardActions, CardContent} from "@mui/material";
import {SearchNodesResourcesTable} from "./SearchNodesTable";
import authProvider from '../../../redux/auth/auth.actions';
import MainListItems from "../../dashboard/listItems";
import {ZeusCopyright} from "../../copyright/ZeusCopyright";
import {resourcesApiGateway} from "../../../gateway/resources";
import {setSearchResources} from "../../../redux/resources/resources.reducer";
import {NodeSearchParams, NodesSlice} from "../../../redux/resources/resources.types";
import {RootState} from "../../../redux/store";

const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
}

export const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

export const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        '& .MuiDrawer-paper': {
            position: 'relative',
            whiteSpace: 'nowrap',
            width: drawerWidth,
            transition: theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
            boxSizing: 'border-box',
            ...(!open && {
                overflowX: 'hidden',
                transition: theme.transitions.create('width', {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.leavingScreen,
                }),
                width: theme.spacing(7),
                [theme.breakpoints.up('sm')]: {
                    width: theme.spacing(9),
                },
            }),
        },
    }),
);

const mdTheme = createTheme();

function SearchComputeDashboardContent() {
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();
    const [loading, setIsLoading] = React.useState(false);
    const resources = useSelector((state: RootState) => state.resources.searchResources);

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    const handleSearchRequest = async () => {
        try {
            setIsLoading(true)
            // setRequestStatus('pending');s
            const CloudProviderRegions: { [key: string]: string[] } = {
                aws: ["us-west-1"],
                do: ["nyc1"],
                gcp: ["us-central1"],
                ovh: ["us-west-or-1"],
            };

            const payloadNodeSearchParams: NodeSearchParams = {
                cloudProviderRegions: CloudProviderRegions,
            };
            const response = await resourcesApiGateway.searchNodeResources(payloadNodeSearchParams);
            if (response.status < 400) {
                const re = response.data as NodesSlice;
                dispatch(setSearchResources(re));
            }

            // if (response.status === 200 || response.status === 202 || response.status === 204) {
            //     setRequestStatus('success');
            //     return
            // } else if (response.status === 403) {
            //     setRequestStatus('missingBilling');
            //     setFreeTrial(true)
            //     return
            // } else if (response.status === 412) {
            //     setRequestStatus('outOfCredits');
            //     setFreeTrial(true)
            //     return
            // } else {
            //     setRequestStatus('error');
            //     return
            // }
        } catch (error: any) {
            // setRequestStatus('error');
            // const status: number = error.response.status;
            // if (status === 403) {
            //     setRequestStatus('missingBilling');
            //     setFreeTrial(true)
            // } else if (status === 412) {
            //     setRequestStatus('outOfCredits');
            //     // Disable the button for 30 seconds
            //     setFreeTrial(true)
            // } else {
            //     setRequestStatus('error');
            // }
        } finally {
            setIsLoading(false)
        }
    };

    if (loading) {
        return <div>Loading...</div>;
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
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Card sx={{ maxWidth: 700 }}>
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="div">
                                    Compute Search Engine
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Search for compute resources by cloud provider, region, slug, and description.
                                </Typography>
                            </CardContent>
                            <CardActions >
                                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                    <Button variant="contained" onClick={handleSearchRequest} >Search</Button>
                                </div>
                            </CardActions>
                        </Card>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <SearchNodesResourcesTable loading={loading} resources={resources}/>
                    </Container>
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function SearchDashboard() {
    return <SearchComputeDashboardContent />;
}
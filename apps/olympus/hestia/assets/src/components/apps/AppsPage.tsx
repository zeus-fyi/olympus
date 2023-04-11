import * as React from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import {AppBar, Drawer} from '../dashboard/Dashboard';
import Box from '@mui/material/Box';
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
import {useDispatch} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import MainListItems from "../dashboard/listItems";
import {PrivateAppsTable} from "./AppsTable";
import {Card, CardContent, Stack} from "@mui/material";
import {PublicAppsTable} from "./PublicAppsTable";

const mdTheme = createTheme();

function AppsPageContent() {
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
                            Apps
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
                    <Container maxWidth="xl" >
                        <div style={{ display: 'flex' }}>
                            <Stack direction="column" spacing={2} sx={{mb: 4 }}>
                                <Container maxWidth="xl" sx={{ mt: 0, mb: 0 }}>
                                    <Card>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Public Registered Apps
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                The table below contains apps template workloads that you can copy, edit, and deploy.
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Container>
                                <Container maxWidth="xl" sx={{ mt: 0, mb: 4 }}>
                                    <PublicAppsTable />
                                </Container>
                            </Stack>
                        </div>
                    </Container>
                    <Container maxWidth="xl" >
                        <div style={{ display: 'flex' }}>
                            <Stack direction="column" spacing={2} sx={{mb: 4 }}>
                                <Container maxWidth="xl" sx={{ mt: 0, mb: 0 }}>
                                    <Card>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Private Registered Apps
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                The table below contains apps that are registered workloads that you can deploy, edit, or upgrade.
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Container>
                                <Container maxWidth="xl" sx={{ mt: 0, mb: 4 }}>
                                    <PrivateAppsTable />
                                </Container>
                            </Stack>
                        </div>
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function AppsPage() {
    return <AppsPageContent />;
}
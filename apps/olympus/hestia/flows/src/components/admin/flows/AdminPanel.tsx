import * as React from 'react';
import {useState} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
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
import {AdminWorkflowViews} from "./AdminWorkflowViews";
import authProvider from "../../../redux/auth/auth.actions";
import {ZeusCopyright} from "../../copyright/ZeusCopyright";
import MainListItems from "../../dashboard/listItems";
import {AppBar, Drawer} from "./AdminHelpers";

const mdTheme = createTheme();

function AdminDashboard() {
    const [open, setOpen] = useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();
    const [loading, setIsLoading] = useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    }
    // admin/dashboard/stats
    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }
    if (loading) {
        return <div>Loading...</div>;
    }
    // getAdminDashboardStats
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
                        <AdminWorkflowViews />
                    </Container>
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function AdminPanel() {
    return <AdminDashboard />;
}
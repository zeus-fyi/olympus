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
import Paper from '@mui/material/Paper';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import {TableContainer, TableRow} from '@mui/material';
import TableBody from "@mui/material/TableBody";
import MainListItems from "../dashboard/listItems";

const mdTheme = createTheme();

function createData(
    cloudCtxNsID: string,
    cloudProvider: string,
    region: string,
    context: string,
    namespace: string,
) {
    return {cloudCtxNsID, cloudProvider, region, context, namespace};
}

const clusterRows = [
    createData('1243535','do', 'sfo3','do-sfo3-zeus', 'eth-indexer'),
    createData('12235535','do', 'sfo3', 'do-sfo3-zeus','ephemeral-staking'),
];

function ClustersContent() {
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
                        {<CloudClusters />}
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function Clusters() {
    return <ClustersContent />;
}

function CloudClusters() {
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>CloudCtxNsID</TableCell>
                        <TableCell align="left">CloudProvider</TableCell>
                        <TableCell align="left">Region</TableCell>
                        <TableCell align="left">Context</TableCell>
                        <TableCell align="left">Namespace</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {clusterRows.map((row) => (
                        <TableRow
                            key={row.cloudCtxNsID}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.cloudCtxNsID}
                            </TableCell>
                            <TableCell align="left">{row.cloudProvider}</TableCell>
                            <TableCell align="left">{row.region}</TableCell>
                            <TableCell align="left">{row.context}</TableCell>
                            <TableCell align="left">{row.namespace}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

const R = require('ramda');

// export const orderMatchFormatMany = R.curry((clusters: any) => R.map(clustersFormat(clusters), clusters));

const mapStateToProps = (state: any) => ({
    clusters: state,
});

const asyncComponentConfig = {
    props: ['clusters'],
    load: async (props: any) => { },
    initial: [],
    dataProp: 'history'
};

// const TradeHistory = connect(mapStateToProps, null)(
//     AsyncComponent(asyncComponentConfig)(TradeHistoryComponent)
// );
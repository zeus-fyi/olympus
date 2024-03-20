import * as React from 'react';
import {useState} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import {AppBar, Drawer} from '../dashboard/Dashboard';
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
import {Card, CardContent, TableContainer, TableRow} from '@mui/material';
import TableBody from "@mui/material/TableBody";
import MainListItems from "../dashboard/listItems";
import {accessApiGateway} from "../../gateway/access";

const mdTheme = createTheme();

function AccessKeys() {
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
                            Access
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
                        {<ApiKeys />}
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function Access() {
    return <AccessKeys />;
}

function createData(
    service: string,
    keyName: string,
    accessKey: string,
) {
    return {service, keyName, accessKey };
}

const initialRows = [
    createData('Zeus', 'API Access Key','*******************************************************************'),
];

function ApiKeys() {
    const [rows, setRows] = useState(initialRows);

    const handleRequestApiKey = async (rowIndex: number) => {
        const response = await accessApiGateway.sendApiKeyGenRequest();
        const data = response.data;
        const updatedRows = rows.map((row, index) => {
            if (index === rowIndex) {
                return { ...row, service: 'Zeus', keyName: data.apiKeyName, accessKey: data.apiKeySecret };
            }
            return row;
        });
        setRows(updatedRows);
    };

    return (
        <div>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
            <Card>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        API Key Access
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        You won't be able to see your API key again after you generate it, so if you lose it you'll have to generate a new one.
                    </Typography>
                </CardContent>
            </Card>
            </Container>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>

                <TableContainer component={Paper}>
            <Table sx={{ minWidth: 300 }} aria-label="simple table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333' }}>
                        <TableCell style={{ color: 'white' }}>Service</TableCell>
                        <TableCell style={{ color: 'white' }} align="left">Key Name</TableCell>
                        <TableCell style={{ color: 'white' }} align="left">Access Key</TableCell>
                        <TableCell style={{ color: 'white' }} align="left"></TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rows.map((row, index) => (
                        <TableRow
                            key={row.service}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.service}
                            </TableCell>
                            <TableCell align="left">{row.keyName}</TableCell>
                            <TableCell align="left" sx={{ maxWidth: 100 }}>{row.accessKey}</TableCell>
                            <TableCell align="right" sx={{ paddingLeft: 1 }}>
                                <Button
                                    color="primary"
                                    variant="contained"
                                    onClick={() => handleRequestApiKey(index)}
                                >
                                    Request API Key
                                </Button>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
            </Container>

        </div>
    );
}
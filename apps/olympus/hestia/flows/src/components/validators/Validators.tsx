import * as React from 'react';
import {useEffect, useState} from 'react';
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
import {Card, CardContent, TableContainer, TableFooter, TablePagination, TableRow} from '@mui/material';
import TableBody from "@mui/material/TableBody";
import MainListItems from "../dashboard/listItems";
import {validatorsApiGateway} from "../../gateway/validators";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";

const mdTheme = createTheme();

function createData(
    network: string,
    groupName: string,
    pubkey: string,
    feeRecipient: string,
    enabled: string,
) {
    return {network, groupName,pubkey,feeRecipient,enabled};
}

function ValidatorsServiceContent() {
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
                            Validators
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
                                    Validator Management for Ethereum Staking
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Currently, you can only view validators which have successfully registered for service. Contact us if you want to register your validator for service.
                                </Typography>
                            </CardContent>
                        </Card>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        {<Validators />}
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function ValidatorsServices() {
    return <ValidatorsServiceContent />;
}

function Validators() {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);

    const [validators, setValidators] = useState([{}]);
    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };
    const handleChangePage = (
        event: React.MouseEvent<HTMLButtonElement> | null,
        newPage: number,
    ) => {
        setPage(newPage);
    };

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - validators.length) : 0;

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await validatorsApiGateway.getValidators();
                const validatorsData: any[] = response.data;
                const validatorRows = validatorsData.map((v: any) =>
                    createData(getNetwork(v.protocolNetworkID), v.groupName, v.pubkey, v.feeRecipient, booleanString(v.enabled))
                );
                setValidators(validatorRows)
            } catch (error) {
                console.log("error", error);
            }}
        fetchData();
    }, []);
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="validators pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Network</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">GroupName</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">PublicKey</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">FeeRecipient</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Enabled</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {(rowsPerPage > 0
                        ? validators.slice(page * rowsPerPage, page*rowsPerPage+rowsPerPage) : validators).map((row: any,i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.network}
                            </TableCell>
                            <TableCell align="left">{row.groupName}</TableCell>
                            <TableCell align="left">{row.pubkey}</TableCell>
                            <TableCell align="left">{row.feeRecipient}</TableCell>
                            <TableCell align="left">{row.enabled}</TableCell>
                        </TableRow>
                    ))}
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={6} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={[10, 25, 100, { label: 'All', value: -1 }]}
                            colSpan={3}
                            count={validators.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            SelectProps={{
                                inputProps: {
                                    'aria-label': 'rows per page',
                                },
                                native: true,
                            }}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            ActionsComponent={TablePaginationActions}
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}

export function getNetwork(networkID: number){
    if (BigInt(networkID) === BigInt(1)) {
        return 'Mainnet'
    }
    if (BigInt(networkID) === BigInt(5)) {
        return 'Goerli'
    }
    if (BigInt(networkID) === BigInt(1673748447294772000)) {
        return 'Ephemery'
    }
    return 'Unknown'
}

export function getNetworkId(network: string){
    if (network === 'Mainnet' || network === 'mainnet') {
        return 1
    }
    if (network === 'Goerli' || network === 'goerli') {
        return 5
    }
    if (network === 'Ephemery' || network === 'ephemery') {
        return  1673748447294772000
    }
    return 1673748447294772000
}

function booleanString(bool: boolean) {
  if (bool) {
      return 'True'
  } else {
    return 'False'
  }
}
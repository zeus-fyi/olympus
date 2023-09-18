import * as React from "react";
import {useEffect} from "react";
import {useDispatch} from "react-redux";
import {useNavigate} from "react-router-dom";
import {appsApiGateway} from "../../gateway/apps";
import {setPrivateOrgApps} from "../../redux/apps/apps.reducer";
import {TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
// @ts-ignore
import {ReactComponent as AvaxLogo} from '../../static/avax-logo.svg';

export function PublicAppsTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const dispatch = useDispatch();
    let navigate = useNavigate();

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await appsApiGateway.getPrivateApps();
                dispatch(setPrivateOrgApps(response));
            } catch (e) {
            }
        }
        fetchData();
    }, [dispatch]);


    const handleClickApp = async (event: any, app: any) => {
        event.preventDefault();
        navigate('/apps/'+app);
    }

    const publicApps = [{category: 'web3', appName: 'Avax', type: 'Cluster'}, {category: 'web3', appName: 'Ethereum', type: 'Cluster'}]

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >App Logo</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >App Category</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Type</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                        <TableRow
                            onClick={(event) => handleClickApp(event, 'microservice')}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell style={{ width: 200 }}>
                                <img src={require("../../static/microservice.png")} alt={'templates'} style={{ width: '25%', height: '40%'}} />
                            </TableCell>
                            <TableCell component="th" scope="row">
                                {'templates'}
                            </TableCell>
                            <TableCell align="left">{'Microservice'}</TableCell>
                            <TableCell align="left">{ 'Cluster'}</TableCell>
                        </TableRow>
                        <TableRow
                            onClick={(event) => handleClickApp(event, 'avax')}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell style={{ width: 200 }}>
                                <AvaxLogo style={{ width: '25%', height: '40%'}} alt='web3' />
                            </TableCell>
                            <TableCell component="th" scope="row">
                                {'web3'}
                            </TableCell>
                            <TableCell align="left">{'Avax'}</TableCell>
                            <TableCell align="left">{ 'Cluster'}</TableCell>
                        </TableRow>
                    <TableRow
                        onClick={(event) => handleClickApp(event, 'eth')}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                        <TableCell style={{ width: 200 }}>
                            <img src={require("../../static/eth.png")} alt={'web3'} style={{ width: '18%', height: '5%', marginLeft: '5px'}} />
                        </TableCell>
                        <TableCell component="th" scope="row">
                            {'web3'}
                        </TableCell>
                        <TableCell align="left">{'Ethereum'}</TableCell>
                        <TableCell align="left">{ 'Cluster'}</TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </TableContainer>
    );
}

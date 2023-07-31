import * as React from "react";
import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {appsApiGateway} from "../../gateway/apps";
import {setPrivateOrgApps} from "../../redux/apps/apps.reducer";
import {RootState} from "../../redux/store";
import {useNavigate} from "react-router-dom";

export function PrivateAppsTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const privateOrgApps = useSelector((state: RootState) => state.apps.privateOrgApps);
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

    const handleClick = async (event: any, app: any) => {
        event.preventDefault();
        navigate('/app/'+app.topologySystemComponentID);
    }

    if (privateOrgApps == null) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - privateOrgApps.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >AppID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Type</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Name</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {privateOrgApps.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            onClick={(event) => handleClick(event, row)}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.topologySystemComponentID}
                            </TableCell>
                            <TableCell align="left">{row.topologyClassTypeID === 4 ? 'Cluster' : 'Matrix'}</TableCell>
                            <TableCell align="left">{row.topologySystemComponentName}</TableCell>
                        </TableRow>
                        ))}
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={4} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={[10, 25, 100, { label: 'All', value: -1 }]}
                            colSpan={4}
                            count={privateOrgApps.length}
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

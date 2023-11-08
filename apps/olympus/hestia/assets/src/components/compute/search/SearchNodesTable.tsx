import * as React from "react";
import {useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../redux/store";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";

export function SearchNodesResourcesTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const resources = useSelector((state: RootState) => state.resources.searchResources);
    const dispatch = useDispatch();
    const [statusMessage, setStatusMessage] = useState('');
    const [statusMessageRowIndex, setStatusMessageRowIndex] = useState<number | null>(null);


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

    const handleRemoveRow = async (rowIndex: number, orgResourceID: number) => {
        try {
           //const response = await resourcesApiGateway.destroyAppResource(orgResourceID);
            setStatusMessageRowIndex(rowIndex);
            setStatusMessage(`OrgResourceID ${orgResourceID} deletion in progress`);
        } catch (error) {
            console.error(error);
            setStatusMessageRowIndex(rowIndex);
            setStatusMessage(`Error deleting resource ID ${orgResourceID}`);
        }
    }
    if (resources === null || resources === undefined) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - resources.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >ResourceID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >CloudProvider</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Region</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Slug</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Hourly Cost</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Monthly Cost</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Description</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {resources.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.resourceID}
                            </TableCell>
                            <TableCell align="left">{row.cloudProvider}</TableCell>
                            <TableCell align="left">{row.region}</TableCell>
                            <TableCell align="left">{row.slug}</TableCell>
                            <TableCell align="left">{(row.priceHourly*1.0).toFixed(2)}</TableCell>
                            <TableCell align="left">{(row.priceMonthly*1.0).toFixed(2)}</TableCell>
                            <TableCell align="left">{row.description}</TableCell>
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
                            count={resources.length}
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

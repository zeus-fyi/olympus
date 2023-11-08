import * as React from "react";
import {useEffect} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";

export function SearchNodesResourcesTable(props: any) {
    const { resources, loading } = props;
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);

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
    useEffect(() => {
    }, [loading]);

    if (loading) {
        return (<div></div>)
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
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >vCPUs</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Memory</TableCell>
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
                            <TableCell align="left">{row.vcpus}</TableCell>
                            <TableCell align="left">{(row.memory / (1024)).toFixed(1) + ' GB'}</TableCell>
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

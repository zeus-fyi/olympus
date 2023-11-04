import * as React from "react";
import {useEffect, useState} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {mevApiGateway} from "../../../gateway/mev";

function createData(
    eventID: string,
    submissionTime: string,
    bundleHash: string,
) {
    return {eventID, submissionTime, bundleHash};
}

export function MevBundlesTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [bundles, setBundles] = useState([{}]);
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
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - bundles.length) : 0;

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await mevApiGateway.getDashboardInfo();
                console.log(response.data)
                const mevDashboardTable: any[] = response.data.bundles;
                const mevDashboardTableRows = mevDashboardTable.map((v: any) =>
                    createData(v.eventID, v.submissionTime, v.bundleHash)
                );
                setBundles(mevDashboardTableRows)
            } catch (error) {
                console.log("error", error);
            }}
        fetchData();
    }, []);
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="mev bundles pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >EventID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Time</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">BundleHash</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {(rowsPerPage > 0
                        ? bundles.slice(page * rowsPerPage, page*rowsPerPage+rowsPerPage) : bundles).map((row: any,i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.eventID}
                            </TableCell>
                            <TableCell align="left">{row.submissionTime}</TableCell>
                            <TableCell align="left">{row.bundleHash}</TableCell>
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
                            count={bundles.length}
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
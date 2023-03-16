import * as React from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";

export function ValidatorsDepositsTable(props: any) {
    const { activeStep, depositData } = props;
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

    if (depositData == null) {
        return (<div></div>)
    }
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - depositData.length) : 0;

    console.log(depositData)
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="validators pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#8991B0'}} >
                        {activeStep === 5 && (
                            <TableCell style={{ fontWeight: 'normal', color: 'white' }} align="left">Verified</TableCell>
                        )}
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Pubkey</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Amount</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Signature</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Withdrawal Credentials</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {(rowsPerPage > 0
                        ? depositData.slice(page * rowsPerPage, page*rowsPerPage+rowsPerPage) : depositData).map((row: any,i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            {activeStep === 5 && (row.pubkey !== undefined) && (
                                <TableCell align="left">{row.verified ? 'True' : 'False'}</TableCell>
                            )}
                            <TableCell component="th" scope="row">
                                {row.pubkey}
                            </TableCell>
                            <TableCell align="left">{row.amount}</TableCell>
                            <TableCell align="left">{row.signature}</TableCell>
                            <TableCell align="left">{row.withdrawal_credentials}</TableCell>
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
                            count={depositData.length}
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
import * as React from "react";
import {useState} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";

function createDepositDataRows(
    pubkey: string,
    feeRecipient: string,
    amount : string,
    signature: string,
    withdrawalCredentials: string,
) {
    return {pubkey, feeRecipient,amount, signature,withdrawalCredentials};
}

export function ValidatorsDepositsTable() {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [validators, setValidators] = useState([{}]);

    const depositsData = useSelector((state: RootState) => state.awsCredentials.depositData);
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

    const depositDataRows = depositsData.map((v: any) =>
        createDepositDataRows(v.pubkey, v.feeRecipient, v.amount, v.signature, v.withdrawal_credentials)
    )

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - depositDataRows.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="validators pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#8991B0'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Pubkey</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Amount</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Signature</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Withdrawal Credentials</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {(rowsPerPage > 0
                        ? depositDataRows.slice(page * rowsPerPage, page*rowsPerPage+rowsPerPage) : depositDataRows).map((row: any,i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
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
                            count={depositDataRows.length}
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
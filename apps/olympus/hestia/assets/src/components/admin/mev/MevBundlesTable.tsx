import * as React from "react";
import {useEffect, useState} from "react";
import {validatorsApiGateway} from "../../../gateway/validators";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";

function MevBundlesTable() {
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
                // const validatorRows = validatorsData.map((v: any) =>
                //     createData(getNetwork(v.protocolNetworkID), v.groupName, v.pubkey, v.feeRecipient, booleanString(v.enabled))
                // );
                // setValidators(validatorRows)
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
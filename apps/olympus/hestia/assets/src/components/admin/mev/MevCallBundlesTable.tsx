import * as React from "react";
import {Box, Collapse, TableContainer, TableFooter, TablePagination, TableRow, Typography} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import {createCallBundleData} from "./Mev";

export function MevCallBundlesTable(props: any) {
    const {callBundles} = props
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
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - callBundles.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="mev call bundles pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left"></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Event Time</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">BundleHash</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Builder</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Bundle GasPrice</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Coinbase Diff</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Gas Fees</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && callBundles && callBundles.map((row: any) => (
                        <CallBundlesRow key={row.eventID} row={row} />
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
                            count={callBundles.length}
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

function CallBundlesRow(props: { row: ReturnType<typeof createCallBundleData> }) {
    const { row } = props;
    const [open, setOpen] = React.useState(false);

    const explorerURL = 'https://etherscan.io';
    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.submissionTime}
                </TableCell>
                <TableCell align="left">{row.bundleHash}</TableCell>
                <TableCell align="left">{row.builderName}</TableCell>
                <TableCell align="left">{row.bundleGasPrice}</TableCell>
                <TableCell align="left">{row.coinbaseDiff}</TableCell>
                <TableCell align="left">{row.gasFees}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Called Bundle Transactions
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>TxHash</TableCell>
                                        <TableCell>GasPrice</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {/*{row && row.results.map((index: number, bundledTxRow: any) => (*/}
                                    {/*    <TableRow key={index}>*/}
                                    {/*        <TableCell component="th" scope="row">*/}
                                    {/*            {bundledTxRow.txHash ? (*/}
                                    {/*                <a*/}
                                    {/*                    href={explorerURL +'/tx/' + bundledTxRow.txHash}*/}
                                    {/*                    target="_blank"*/}
                                    {/*                    rel="noopener noreferrer"*/}
                                    {/*                >*/}
                                    {/*                    {bundledTxRow.txHash.slice(0, 40)}*/}
                                    {/*                </a>*/}
                                    {/*            ) : 'None'}*/}
                                    {/*        </TableCell>*/}
                                    {/*        <TableCell>*/}
                                    {/*            {bundledTxRow.submissionTime}*/}
                                    {/*        </TableCell>*/}
                                    {/*        <TableCell>*/}
                                    {/*            {bundledTxRow.builderHash}*/}
                                    {/*        </TableCell>*/}
                                    {/*        <TableCell>*/}
                                    {/*            {bundledTxRow.bundleGasPrice}*/}
                                    {/*        </TableCell>*/}
                                    {/*        <TableCell>*/}
                                    {/*            {bundledTxRow.coinbaseDiff}*/}
                                    {/*        </TableCell>*/}
                                    {/*        /!*<TableCell>{(bundledTxRow.ethTxGas.gasPrice.Int64 / 1e9).toLocaleString('fullwide', { useGrouping: false })}</TableCell>*!/*/}
                                    {/*    </TableRow>*/}
                                    {/*))}*/}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}
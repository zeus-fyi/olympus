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
import {createBundleData} from "./Mev";

export function MevBundlesTable(props: any) {
    const {bundles} =props
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
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - bundles.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="mev bundles pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left"></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Event Time</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">BundleHash</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && bundles && bundles.map((row: any) => (
                        <Row key={row.name} row={row} />
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

function Row(props: { row: ReturnType<typeof createBundleData> }) {
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
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Bundled Transactions
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>TxHash</TableCell>
                                        <TableCell>Status</TableCell>
                                        <TableCell>BlockNumber</TableCell>
                                        <TableCell>TxIndex</TableCell>
                                        <TableCell>TxFee</TableCell>
                                        <TableCell>GasUsed</TableCell>
                                        <TableCell>EffectiveGasPrice</TableCell>
                                        <TableCell>GasFeeCap</TableCell>
                                        <TableCell>GasTipCap</TableCell>
                                        <TableCell>GasLimit</TableCell>
                                        {/*<TableCell>GasPrice</TableCell>*/}
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.bundledTxs.map((bundledTxRow) => (
                                        <TableRow key={bundledTxRow.ethTx.eventID}>
                                            <TableCell component="th" scope="row">
                                                {bundledTxRow.ethTx.txHash ? (
                                                    <a
                                                        href={explorerURL +'/tx/' + bundledTxRow.ethTx.txHash}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                    >
                                                        {bundledTxRow.ethTx.txHash.slice(0, 40)}
                                                    </a>
                                                ) : 'None'}
                                            </TableCell>
                                            <TableCell>{bundledTxRow.ethTxReceipts.status}</TableCell>
                                            <TableCell>{bundledTxRow.ethTxReceipts.blockNumber}</TableCell>
                                            <TableCell>{bundledTxRow.ethTxReceipts.transactionIndex}</TableCell>
                                            <TableCell>{(((bundledTxRow.ethTxReceipts.effectiveGasPrice / 1e18))*(bundledTxRow.ethTxReceipts.gasUsed)).toFixed(5)}</TableCell>
                                            <TableCell>{bundledTxRow.ethTxReceipts.gasUsed}</TableCell>
                                            <TableCell>{(bundledTxRow.ethTxReceipts.effectiveGasPrice / 1e9).toLocaleString('fullwide', { useGrouping: false })}</TableCell>
                                            <TableCell>{(bundledTxRow.ethTxGas.gasFeeCap.Int64 / 1e9).toLocaleString('fullwide', { useGrouping: false })}</TableCell>
                                            <TableCell>{(bundledTxRow.ethTxGas.gasTipCap.Int64 / 1e9).toLocaleString('fullwide', { useGrouping: false })}</TableCell>
                                            <TableCell>{bundledTxRow.ethTxGas.gasLimit.Int64}</TableCell>
                                            {/*<TableCell>{(bundledTxRow.ethTxGas.gasPrice.Int64 / 1e9).toLocaleString('fullwide', { useGrouping: false })}</TableCell>*/}
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell align="left">User</TableCell>
                                        <TableCell align="left">Total Tx Fees</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {Object.entries(row.traderInfo).map(([traderKey, info]) => (
                                        <TableRow key={traderKey}>
                                            <TableCell component="th" scope="row" align="left">
                                                    <a
                                                        href={explorerURL +'/address/' + traderKey}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                    >
                                                        {traderKey.slice(0, 40)}
                                                    </a>
                                            </TableCell>
                                            <TableCell align="left">
                                                {info.totalTxFees.toFixed(5)} {'Eth'}{/* Displaying the total tx fees in eth*/}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}


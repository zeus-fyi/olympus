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
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Actual Profit</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Worst Case Profit</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Bundle GasPrice</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Coinbase Diff</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Gas Fees</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} align="left">Status</TableCell>
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
                <TableCell align="left">{row.bundleHash.slice(0,40)}</TableCell>
                <TableCell align="left">{row.builderName}</TableCell>
                <TableCell align="left">{row.actualProfitAmountOut + ' Eth'}</TableCell>
                <TableCell align="left">{row.expectedProfitAmountOut + ' Eth'}</TableCell>
                <TableCell align="left">{row.bundleGasPrice + ' Gwei'}</TableCell>
                <TableCell align="left">{row.coinbaseDiff + ' Eth'}</TableCell>
                <TableCell align="left">{row.gasFees + ' Eth'}</TableCell>
                <TableCell align="left">{row.status}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Called Bundle Transactions: {row.tradeMethod}
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>TxHash</TableCell>
                                        <TableCell>Coinbase Diff</TableCell>
                                        <TableCell>GasFees</TableCell>
                                        <TableCell>GasPrice</TableCell>
                                        <TableCell>Error</TableCell>
                                        <TableCell>Revert</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                       {row.results.map((txRow) => (
                                            <TableRow key={txRow.txHash}>
                                                <TableCell component="th" scope="row">
                                                    <a href={`${explorerURL}/tx/${txRow.txHash}`} target="_blank" rel="noreferrer">{txRow.txHash}</a>
                                                </TableCell>
                                                <TableCell>{txRow.coinbaseDiff}</TableCell>
                                                <TableCell>{txRow.gasFees}</TableCell>
                                                <TableCell>{txRow.gasPrice}</TableCell>
                                                <TableCell>{txRow.error}</TableCell>
                                                <TableCell>{txRow.revert}</TableCell>
                                            </TableRow>
                                        ))}
                                </TableBody>
                            </Table>
                        </Box>
                        <Box sx={{ margin: 2 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Block Analysis
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Pair Address</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow >
                                        <TableCell component="th" scope="row">
                                            <a href={`${explorerURL}/address/${row.pairAddress}`} target="_blank" rel="noreferrer">{row.pairAddress}</a>
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Box>
                        <Box sx={{ margin: 2 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Trade Analysis
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Block Number Seen</TableCell>
                                        <TableCell>Block Number Confirmed</TableCell>
                                        <TableCell>Tx Index</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow >
                                        <TableCell component="th" scope="row">
                                            <a href={`${explorerURL}/txs?block=${row.seenAtBlockNumber}`} target="_blank" rel="noreferrer">{row.seenAtBlockNumber}</a>
                                        </TableCell>
                                        <TableCell component="th" scope="row">
                                            <a href={`${explorerURL}/txs?block=${row.blockNumber}`} target="_blank" rel="noreferrer">{row.blockNumber}</a>
                                        </TableCell>
                                        <TableCell>{row.transactionIndex}</TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Box>
                        <Box sx={{ margin: 2 }}>
                            <Table size="small" aria-label="purchases">
                                <TableBody >
                                    <TableRow>
                                        <TableCell>Amount In</TableCell>
                                        <TableCell>In Addr</TableCell>
                                        <TableCell>Amount Out</TableCell>
                                        <TableCell>Out Addr</TableCell>
                                    </TableRow>
                                    {row.trades.map((trade: any, ind: number) => (
                                        <TableRow key={ind}>
                                            <TableCell>{trade.amountIn}</TableCell>
                                            <TableCell>
                                                <a href={`${explorerURL}/address/${trade.amountInAddr}`} target="_blank" rel="noreferrer">{trade.amountInAddr}</a>
                                            </TableCell>
                                            <TableCell>{trade.amountOut}</TableCell>
                                            <TableCell>
                                                <a href={`${explorerURL}/address/${trade.amountOutAddr}`} target="_blank" rel="noreferrer">{trade.amountOutAddr}</a>
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
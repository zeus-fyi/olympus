import * as React from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Checkbox from "@mui/material/Checkbox";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import TextField from "@mui/material/TextField";

export function LoadBalancingRoutesTable(props: any) {
    const { loading,rowsPerPage, page,selected, endpoints, handleSelectAllClick, handleClick,
        handleChangeRowsPerPage,handleChangePage, groupName,
        isAdding, setIsAdding, newEndpoint, setNewEndpoint, isUpdatingGroup, handleAddGroupTableEndpointsSubmission,
        handleSubmitNewEndpointSubmission, handleDeleteEndpointsSubmission, handleUpdateGroupTableEndpointsSubmission
    } = props

    if (loading == true) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }
    let safeEndpoints = endpoints ?? [];
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - safeEndpoints.length) : 0;

    return (
        <div>
            <Box sx={{ mt: 4, mb: 4 }}>
                {selected.length > 0 && !isUpdatingGroup && (
                    <Box sx={{ mb: 2 }}>
                        <span>({selected.length} selected endpoints)</span>
                        <Button variant="outlined" color="secondary" onClick={handleDeleteEndpointsSubmission} style={{marginLeft: '10px'}}>
                            Delete
                        </Button>
                    </Box>
                )}
                {selected.length > 0 && isUpdatingGroup && (
                    <Box sx={{ mb: 2 }}>
                        <span>({selected.length} selected endpoints)</span>
                        <Button variant="outlined" color="secondary" onClick={handleAddGroupTableEndpointsSubmission} style={{marginLeft: '10px'}}>
                            Add
                        </Button>
                    </Box>
                )}
                    <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 650 }} aria-label="simple table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333'}} >
                                <TableCell padding="checkbox">
                                    <Checkbox
                                        color="primary"
                                        indeterminate={safeEndpoints.length > 0 && selected.length < safeEndpoints.length}
                                        checked={safeEndpoints.length > 0 && selected.length === safeEndpoints.length}
                                        onChange={handleSelectAllClick}
                                    />
                                </TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Endpoint</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {isAdding && (
                                <TableRow>
                                    <TableCell padding="checkbox"></TableCell>
                                    <TableCell component="th" scope="row">
                                        <Box display="flex" alignItems="center" gap={2}>
                                            <TextField
                                                value={newEndpoint}
                                                onChange={event => setNewEndpoint(event.target.value)}
                                                sx={{ height: '53px', flex: 1 }} // adjust height here
                                            />
                                            <Button
                                                variant="contained"
                                                color="primary"
                                                onClick={handleSubmitNewEndpointSubmission}
                                            >
                                                Submit
                                            </Button>
                                        </Box>
                                    </TableCell>
                                </TableRow>
                            )}
                            {safeEndpoints.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row: string, i: number) => (
                                <TableRow
                                    key={i}
                                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                >
                                    <TableCell padding="checkbox">
                                        <Checkbox
                                            checked={selected.indexOf(row) !== -1}
                                            onChange={() => handleClick(row)}
                                            color="primary"
                                        />
                                    </TableCell>
                                    <TableCell component="th" scope="row">
                                        {row}
                                    </TableCell>
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
                                    count={safeEndpoints.length}
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
            </Box>
        </div>
    );
}

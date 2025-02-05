import * as React from "react";
import {useEffect} from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../redux/store";
import {loadBalancingApiGateway} from "../../../gateway/loadbalancing";
import {setTableMetrics} from "../../../redux/loadbalancing/loadbalancing.reducer";

export function LoadBalancingPriorityScoreMetricsTable(props: any) {
    const { loading, rowsPerPage, page,selected, endpoints, handleSelectAllClick, handleClick,
        handleChangeRowsPerPage,handleChangePage, groupName, selectedTab, handleTabChange,
        isAdding, setIsAdding, newEndpoint, setNewEndpoint, isUpdatingGroup, handleAddGroupTableEndpointsSubmission,
        handleSubmitNewEndpointSubmission, handleDeleteEndpointsSubmission, handleUpdateGroupTableEndpointsSubmission
    } = props
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    const dispatch = useDispatch();
    const [loadingMetrics, setLoadinMetrics] = React.useState(false);
    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoadinMetrics(true); // Set loading to true
                const response = await loadBalancingApiGateway.getTableMetrics(groupName);
                dispatch(setTableMetrics(response.data));
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoadinMetrics(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData();
    }, [groupName]);

    if (loadingMetrics) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }
    if (tableMetrics == null || tableMetrics.metrics == null ||  Object.keys(tableMetrics.metrics).length == 0) {
        return <div></div>
    }
    let safeEndpoints = tableMetrics.routes ?? [];
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
                    <Table sx={{ minWidth: 650 }} aria-label="metric slice table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333'}}>
                                <TableCell style={{ color: 'white'}} align="left">Endpoint</TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Score</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {safeEndpoints && safeEndpoints.map((slice, index) => (
                                <TableRow key={index}>
                                    <TableCell align="left">{slice.Member}</TableCell>
                                    <TableCell align="left">{slice.Score}</TableCell>
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

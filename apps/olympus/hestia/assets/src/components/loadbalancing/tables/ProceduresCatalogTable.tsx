import * as React from "react";
import {useEffect} from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import {Stack, TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../redux/store";
import {loadBalancingApiGateway} from "../../../gateway/loadbalancing";
import {setProceduresCatalog} from "../../../redux/loadbalancing/loadbalancing.reducer";
import ExamplePageMarkdownText from "../markdown/ExamplePageMarkdown";
import {
    avaxMaxBlockAggReduceExample,
    btcMaxBlockAggReduceExample,
    ethMaxBlockAggReduceExample,
    nearMaxBlockAggReduceExample
} from "../markdown/ExampleRequests";
import {IrisApiGateway} from "../../../gateway/iris";

export function ProceduresCatalogTable(props: any) {
    const { loading, rowsPerPage, page,selected, endpoints, handleSelectAllClick, handleClick,
        handleChangeRowsPerPage,handleChangePage, groupName, selectedTab, handleTabChange, selectedMainTab,
        isAdding, setIsAdding, newEndpoint, setNewEndpoint, isUpdatingGroup, handleAddGroupTableEndpointsSubmission,
        handleSubmitNewEndpointSubmission, handleDeleteEndpointsSubmission, handleUpdateGroupTableEndpointsSubmission
    } = props
    const proceduresCatalog = useSelector((state: RootState) => state.loadBalancing.proceduresCatalog);
    const dispatch = useDispatch();
    const [loadingProcedures, setLoadingProcedures] = React.useState(false);
    const [code, setCode] = React.useState('');
    const [showDetails, setShowDetails] = React.useState(false);
    const [showDetailsRow, setShowDetailsRow] = React.useState(-1);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // if (selectedMainTab !== 1 && selectedTab !== 4) {
                //     return
                // }
                setLoadingProcedures(true); // Set loading to true
                const response = await loadBalancingApiGateway.getProceduresCatalog();
                dispatch(setProceduresCatalog(response.data));
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoadingProcedures(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData();
    }, [selectedMainTab, selectedTab]);

    if (loadingProcedures) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }
    if (proceduresCatalog == null || proceduresCatalog.length === 0) {
        return <div></div>
    }
    let safeEndpoints = proceduresCatalog ?? [];

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - safeEndpoints.length) : 0;

    const onChange = async (textInput: string) => {
        setCode(textInput);
    };

    const onViewDetails = async (index: number, procName: string) => {
        if (index === showDetailsRow) {
            setShowDetailsRow(-1);
            setShowDetails(false);
            setCode('')
            return;
        } else {
            setShowDetails(true);
        }

        switch (procName) {
            case 'eth_maxBlockAggReduce':
                setCode(ethMaxBlockAggReduceExample);
                break;
            case 'avax_maxBlockAggReduce':
                setCode(avaxMaxBlockAggReduceExample);
                break;
            case 'near_maxBlockAggReduce':
                setCode(nearMaxBlockAggReduceExample);
                break;
            case 'btc_maxBlockAggReduce':
                setCode(btcMaxBlockAggReduceExample);
                break;
            default:
                break;
        }
        setShowDetailsRow(index);
    };

    const onSubmitPayload = async () => {
        try {
            setLoadingProcedures(true); // Set loading to true
            const response = await IrisApiGateway.sendJsonRpcRequest(groupName, code);
            console.log("response", response.data)
            setCode(response.data);
        } catch (error) {
            console.log("error", error);
        } finally {
            setLoadingProcedures(false); // Set loading to false regardless of success or failure.
        }
    };

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
                                <TableCell style={{ color: 'white'}} align="left">Name</TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Protocol</TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Description</TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {safeEndpoints && safeEndpoints.map((slice, index) => (
                                <TableRow key={index}>
                                    <TableCell align="left">{slice.name}</TableCell>
                                    <TableCell align="left">{slice.protocol}</TableCell>
                                    <TableCell align="left">{slice.description}</TableCell>
                                    <TableCell align="left">
                                        <Stack direction={"row"} spacing={2}>
                                            <Box sx={{ mt: 4, mb: 4 }}>
                                                {
                                                   showDetails && showDetailsRow === index ? (
                                                       <Button variant="contained" color="primary"  onClick={() => onViewDetails(index, '')}>
                                                           Hide Details
                                                       </Button>
                                                    ) : (
                                                       <Button variant="contained" color="primary"  onClick={() => onViewDetails(index, slice.name)}>
                                                           View Details
                                                       </Button>
                                                    )
                                                }
                                            </Box>
                                        </Stack>
                                        {/*<Box sx={{ mt: 4, mb: 4 }}>*/}
                                        {/*    <Button variant="contained" color="primary"  onClick={() => ({})}>*/}
                                        {/*        Settings*/}
                                        {/*    </Button>*/}
                                        {/*</Box>*/}
                                        {/*{statusMessageRowIndex === i && <div>{statusMessage}</div>}*/}
                                    </TableCell>
                                    <TableCell align="left">
                                        {
                                            showDetails && showDetailsRow === index && (
                                                <div>
                                                    <Stack direction={"column"} spacing={2}>
                                                        <ExamplePageMarkdownText onChange={onChange} code={code} setCode={setCode}/>
                                                        {(groupName !== "-all" && groupName !== "unused") && (
                                                            <Box sx={{ mt: 4, mb: 4 }}>
                                                                <Button variant="contained" fullWidth={true} color="primary"  onClick={onSubmitPayload}>
                                                                    Send Request
                                                                </Button>
                                                            </Box>
                                                        )}
                                                    </Stack>
                                                </div>
                                            )
                                        }
                                        {/*{statusMessageRowIndex === i && <div>{statusMessage}</div>}*/}
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

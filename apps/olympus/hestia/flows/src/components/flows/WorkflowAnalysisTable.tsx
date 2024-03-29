import * as React from "react";
import {useEffect} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {aiApiGateway} from "../../gateway/ai";
import {
    setAssistants,
    setEvalFns,
    setOpenRunsRow,
    setOrchDetails,
    setRetrievals,
    setRuns,
    setSchemas,
    setSelectedRuns,
    setTriggerActions,
    setWorkflows
} from "../../redux/ai/ai.reducer";
import {useDispatch, useSelector} from "react-redux";
import Checkbox from "@mui/material/Checkbox";
import {OrchestrationsAnalysis} from "../../redux/ai/ai.types.runs";
import {WorkflowAnalysisRow} from "./WorkflowAnalysisRow";

export function WorkflowAnalysisTable(props: any) {
    const { csvExport } = props;
    const [page, setPage] = React.useState(0);
    const openRunsRow = useSelector((state: any) => state.ai.openRunsRow);
    const orchDetails = useSelector((state: any) => state.ai.orchDetails);
    const selectedRuns = useSelector((state: any) => state.ai.selectedRuns);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
    const workflows = useSelector((state: any) => state.ai.runs);
    const dispatch = useDispatch();

    const fetchRun = async (index: number) => {
        const runId = workflows[index].orchestration.orchestrationStrID;
        try {
            setIsLoading(true);
            const response = await aiApiGateway.getRun(runId);
            // Assuming response.data is an array of OrchestrationsAnalysis
            console.log("response", response.data)
            const runToUpdate: OrchestrationsAnalysis[] = response.data.filter((run: OrchestrationsAnalysis) => run.orchestration.orchestrationStrID === runId);
            if (runToUpdate.length > 0) {
                // Assuming we want to update the orchDetails state with these details
                dispatch(setOrchDetails({ [runId]: runToUpdate[0] }));  // Update orchDetails with the first item matching the criteria
            }
        } catch (error) {
            console.log("error", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleOpen = (index: number) => {
        const isOpen = openRunsRow[index] || false;
        if (!isOpen) {
            fetchRun(index);
        }
        dispatch(setOpenRunsRow({ ...openRunsRow, [index]: !isOpen }));
    };

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

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                setIsLoading(true); // Set loading to true
                const response = await aiApiGateway.getRunsUI();
                dispatch(setRuns(response.data));
            } catch (error) {
                console.log("error", error);
            } finally {
                setIsLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData({});
    }, []);

    if (loading) {
        return (<div></div>)
    }
    if (workflows === null || workflows === undefined) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - workflows.length) : 0;

    const handleClick = (name: string) => {
        const currentIndex = selectedRuns.indexOf(name);
        const newSelected = [...selectedRuns];

        if (currentIndex === -1) {
            newSelected.push(name);
        } else {
            newSelected.splice(currentIndex, 1);
        }
        dispatch(setSelectedRuns(newSelected));
    };

    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.checked) {
            const newSelected = workflows.map((wf: any, ind: number) => ind);
            dispatch(setSelectedRuns(newSelected));
            return;
        }
        dispatch(setSelectedRuns([]));
    };
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell padding="checkbox">
                            <Checkbox
                                color="primary"
                                indeterminate={workflows.length > 0 && selectedRuns.length < workflows.length && selectedRuns.length > 0}
                                checked={workflows.length > 0 && selectedRuns.length === workflows.length}
                                onChange={handleSelectAllClick}
                            />
                        </TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} ></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Run ID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Group</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Type</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Active</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Run Cycles</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Tokens Used</TableCell>
                        {csvExport && <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Export</TableCell>
                        }
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && workflows && workflows.map((row: OrchestrationsAnalysis, index: number) => {
                        // Use orchestrationStrID to find details in orchDetails, or fallback to row if not found
                        const detailedRow = orchDetails[row.orchestration.orchestrationStrID] || row;

                        return (
                            <WorkflowAnalysisRow
                                key={index}
                                row={detailedRow} // Pass detailedRow, which may contain additional details
                                open={openRunsRow[index] || false}
                                handleOpen={() => handleOpen(index)}
                                index={index}
                                csvExport={csvExport}
                                handleClick={handleClick}
                                checked={selectedRuns.indexOf(index) >= 0}
                            />
                        );
                    })}
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
                            count={workflows.length}
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

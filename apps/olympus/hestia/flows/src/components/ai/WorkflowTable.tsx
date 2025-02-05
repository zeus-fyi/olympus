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
    setRetrievals,
    setSchemas,
    setSelectedWorkflows,
    setTriggerActions,
    setWorkflows
} from "../../redux/ai/ai.reducer";
import {useDispatch, useSelector} from "react-redux";
import Checkbox from "@mui/material/Checkbox";
import {Task, WorkflowTemplate} from "../../redux/ai/ai.types";
import {WorkflowRow} from "./WorkflowRow";

export function WorkflowTable(props: any) {
    const {csvFilter} = props;
    const [page, setPage] = React.useState(0);
    const selected = useSelector((state: any) => state.ai.selectedWorkflows);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
    let workflows = useSelector((state: any) => state.ai.workflows);
    if (csvFilter && workflows) {
        workflows = workflows.filter((workflow: WorkflowTemplate) => {
            return workflow.tasks.some((task: Task) => task.responseFormat === 'csv');
        })
    }
    const dispatch = useDispatch();
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
        const fetchData = async () => {
            try {
                setIsLoading(true)
                const response = await aiApiGateway.getWorkflowsRequest();
                const statusCode = response.status;
                if (statusCode < 400) {
                    const data = response.data;
                    dispatch(setWorkflows(data.workflows));
                    dispatch(setEvalFns(data.evalFns));
                    dispatch(setTriggerActions(data.triggerActions));
                    // dispatch(setAiTasks(data.tasks));
                    dispatch(setRetrievals(data.retrievals));
                    dispatch(setAssistants(data.assistants));
                    dispatch(setSchemas(data.schemas))
                } else {
                    console.log('Failed to get workflows', response);
                }
            } catch (e) {
            } finally {
                setIsLoading(false);
            }
        }
        fetchData();
    }, []);

    if (loading) {
        return (<div></div>)
    }
    if (workflows === null || workflows === undefined) {
        return (<div></div>)
    }
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - workflows.length) : 0;

    const handleClick = (index: number) => {
        const currentIndex = selected.indexOf(index);
        const newSelected = [...selected];
        if (currentIndex === -1) {
            newSelected.push(index);
        } else {
            newSelected.splice(currentIndex, 1);
        }
        dispatch(setSelectedWorkflows(newSelected));
    };
    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.checked) {
            const newSelected = workflows.map((wf: any, index: number) => index);
            dispatch(setSelectedWorkflows(newSelected));
            return;
        }
        dispatch(setSelectedWorkflows([] as string[]));
    };
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell padding="checkbox">
                            <Checkbox
                                color="primary"
                                indeterminate={workflows.length > 0 && selected.length < workflows.length && selected.length > 0}
                                checked={workflows.length > 0 && selected.length === workflows.length}
                                onChange={handleSelectAllClick}
                            />
                        </TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} ></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Workflow ID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Group</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Base Period</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && workflows && workflows.map((row: WorkflowTemplate, index: number) => (
                        <WorkflowRow
                            key={index}
                            row={row}
                            index={index}
                            handleClick={() =>handleClick(index)}
                            checked={selected.indexOf(index) >= 0 || false}
                        />
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

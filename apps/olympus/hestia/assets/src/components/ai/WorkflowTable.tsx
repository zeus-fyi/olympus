import * as React from "react";
import {useEffect, useState} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {aiApiGateway} from "../../gateway/ai";
import {setAiTasks, setRetrievals, setWorkflows} from "../../redux/ai/ai.reducer";
import {useDispatch, useSelector} from "react-redux";
import Checkbox from "@mui/material/Checkbox";

export function WorkflowTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [selected, setSelected] = useState<string[]>([]);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
    const workflows = useSelector((state: any) => state.ai.workflows);
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
                    dispatch(setAiTasks(data.tasks));
                    dispatch(setRetrievals(data.retrievals));
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

    const handleClick = (name: string) => {
        const currentIndex = selected.indexOf(name);
        const newSelected = [...selected];

        if (currentIndex === -1) {
            newSelected.push(name);
        } else {
            newSelected.splice(currentIndex, 1);
        }

        setSelected(newSelected);
    };
    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.checked) {
            const newSelected = workflows.map((wf: any) => wf);
            setSelected(newSelected);
            return;
        }
        setSelected([]);
    };
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell padding="checkbox">
                            <Checkbox
                                color="primary"
                                indeterminate={workflows.length > 0 && selected.length < workflows.length}
                                checked={workflows.length > 0 && selected.length === workflows.length}
                                onChange={handleSelectAllClick}
                            />
                        </TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Workflow ID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Group</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Base Period</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {workflows.map((row: any, i: number) => (
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
                                {row.workflowID}
                            </TableCell>
                            <TableCell component="th" scope="row">
                                {row.workflowName}
                            </TableCell>
                            <TableCell component="th" scope="row">
                                {row.workflowGroup}
                            </TableCell>
                            <TableCell component="th" scope="row">
                                {row.fundamentalPeriod + ' ' + row.fundamentalPeriodTimeUnit}
                            </TableCell>
                            {/*<TableCell component="th" scope="row">*/}
                            {/*    {row.active ? 'Yes' : 'No'}*/}
                            {/*</TableCell>*/}
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

import * as React from "react";
import {useState} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useDispatch, useSelector} from "react-redux";
import {TasksRow} from "./TasksRow";

export function TasksTable(props: any) {
    const {taskType} = props;
    const [page, setPage] = React.useState(0);
    const [selected, setSelected] = useState<string[]>([]);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
    const allTasks = useSelector((state: any) => state.ai.tasks);
    const tasks = allTasks.filter((task: any) => task.taskType === taskType);

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

    if (loading) {
        return (<div></div>)
    }
    if (tasks === null || tasks === undefined) {
        return (<div></div>)
    }
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - tasks.length) : 0;

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
            const newSelected = tasks.map((wf: any) => wf);
            setSelected(newSelected);
            return;
        }
        setSelected([]);
    };

    if (emptyRows) {
        return (<div></div>)
    }
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        {/*<TableCell padding="checkbox">*/}
                        {/*    <Checkbox*/}
                        {/*        color="primary"*/}
                        {/*        indeterminate={tasks.length > 0 && selected.length < tasks.length}*/}
                        {/*        checked={tasks.length > 0 && selected.length === tasks.length}*/}
                        {/*        onChange={handleSelectAllClick}*/}
                        {/*    />*/}
                        {/*</TableCell>*/}
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} ></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Task ID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Group</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Model</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && tasks && tasks.map((row: any) => (
                        <TasksRow key={tasks.taskID} row={row} />
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
                            colSpan={4}
                            count={tasks.length}
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


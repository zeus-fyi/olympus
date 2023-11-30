import * as React from "react";
import {Checkbox, TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useDispatch} from "react-redux";
import {TasksRow} from "./TasksRow";

export function TasksTable(props: any) {
    const {selected, tasks, handleClick,handleSelectAllClick} = props;
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
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

    if (tasks === null || tasks === undefined) {
        return (<div></div>)
    }
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - tasks.length) : 0;


    if (emptyRows) {
        return (<div></div>)
    }
    const countTrueValues = (obj: { [key: number]: boolean }): number => {
        return Object.values(obj).filter(value => value).length;
    };
    if (loading) {
        return (<div></div>)
    }

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell padding="checkbox">
                            <Checkbox
                                color="primary"
                                indeterminate={countTrueValues(selected) > 0 && countTrueValues(selected) < tasks.length}
                                checked={(countTrueValues(selected) === tasks.length) && tasks.length > 0}
                                onChange={handleSelectAllClick}
                            />
                        </TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} ></TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Task ID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Group</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Name</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Model</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && tasks && tasks.map((row: any, index: number) => (
                        <TasksRow
                            key={tasks.taskID}
                            row={row}
                            index={index}
                            handleClick={handleClick}
                            checked={selected[index] || false}
                        />
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


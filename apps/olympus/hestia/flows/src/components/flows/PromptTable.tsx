import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import * as React from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

export function TaskPromptsTable(props: any) {
    const {} = props;
    const prompts = useSelector((state: RootState) => state.flows.uploadContentTasks);
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const csvHeaders = prompts.length > 0 ? Object.keys(prompts[0]) : [];

    const handleChangePage = (event: React.MouseEvent<HTMLButtonElement> | null, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0); // Reset to the first page
    };

    const emptyRows = page > 0 ? Math.max(0, (1 + page) * rowsPerPage - prompts.length) : 0;
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 500 }} aria-label="prompts table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333' }}>
                        {csvHeaders.map((header) => (
                            <TableCell key={header} style={{ fontWeight: 'bold', color: 'white' }} align="left">
                                {header.toUpperCase()}
                            </TableCell>
                        ))}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {
                        Array.isArray(prompts) && prompts.length > 0
                            ? prompts
                                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                .map((row: any, index: number) => (
                                    <TableRow key={index} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                                        {csvHeaders.map((header) => (
                                            <TableCell key={`${index}-${header}`} align="left">
                                                {row[header]}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                ))
                            : <div></div> // Or render a placeholder row/message here if contacts is empty or not an array
                    }
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={csvHeaders.length} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={[10, 25, 100, { label: 'All', value: -1 }]}
                            colSpan={csvHeaders.length}
                            count={prompts.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            component="div" // Ensure correct rendering of pagination
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}
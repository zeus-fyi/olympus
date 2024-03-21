import * as React from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";

export function ContactsTable(props: any) {
    const {} = props;
    const contacts = useSelector((state: RootState) => state.flows.uploadContentContacts);
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    // Extract CSV headers if contacts is not empty
    const csvHeaders = contacts.length > 0 ? Object.keys(contacts[0]) : [];

    const handleChangePage = (event: React.MouseEvent<HTMLButtonElement> | null, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0); // Reset to the first page
    };

    const emptyRows = page > 0 ? Math.max(0, (1 + page) * rowsPerPage - contacts.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 500 }} aria-label="contacts table">
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
                        // Check if contacts is an array and has length before attempting to slice and map
                        Array.isArray(contacts) && contacts.length > 0
                            ? contacts
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
                            count={contacts.length}
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
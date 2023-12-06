import * as React from "react";
import {Collapse, TableBody, TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Table from "@mui/material/Table";

export function RetrievalsRow(props: { row: ReturnType<typeof createRetrievalDetailsData>, index: number, handleClick: any, checked: boolean}) {
    const { row, index, handleClick, checked } = props;
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell align="center" >
                    <Checkbox
                        checked={checked}
                        onChange={() => handleClick(index)}
                        color="primary"
                    />
                </TableCell>
                <TableCell align="left">{row.retrievalID}</TableCell>
                <TableCell align="left">{row.retrievalGroup}</TableCell>
                <TableCell align="left">{row.retrievalName}</TableCell>
                <TableCell align="left">{row.retrievalPlatform}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Retrieval Details
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableBody>
                                    <TableRow >
                                        <TableCell component="th" scope="row">
                                            {prettyPrintJSON(row.instructions)}
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}

export const prettyPrintJSON = (byteArrayString: string): string => {
    try {
        // Assuming byteArrayString is a base64 encoded string of the byte array
        const decodedString = atob(byteArrayString);
        const jsonObject = JSON.parse(decodedString);
        return JSON.stringify(jsonObject, null, 2); // Pretty print with 2 spaces indentation
    } catch (error) {
        console.error('Error parsing JSON:', error);
        return byteArrayString; // Fallback to original string in case of error
    }
};
export function createRetrievalDetailsData(
    retrievalID: number,
    retrievalName: string,
    retrievalGroup: string = 'default',
    retrievalPlatform: string,
    instructions: string = '',

) {
    return {
        retrievalID,
        retrievalName,
        retrievalGroup,
        retrievalPlatform,
        instructions,
    };
}

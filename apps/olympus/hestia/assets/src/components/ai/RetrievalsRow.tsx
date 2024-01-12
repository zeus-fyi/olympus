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
import Button from "@mui/material/Button";
import {setRetrieval} from "../../redux/ai/ai.reducer";
import {Retrieval} from "../../redux/ai/ai.types2";
import {useDispatch} from "react-redux";

export function RetrievalsRow(props: { row: Retrieval, index: number, handleClick: any, checked: boolean}) {
    const { row, index, handleClick, checked } = props;
    const [open, setOpen] = React.useState(false);
    const dispatch = useDispatch();
    const handleEditRetrieval = async (e: any, ret: Retrieval) => {
        dispatch(setRetrieval(ret))
    }
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
                <TableCell align="left">{row.retrievalID? row.retrievalID : 0}</TableCell>
                <TableCell align="left">{row.retrievalGroup}</TableCell>
                <TableCell align="left">{row.retrievalName}</TableCell>
                <TableCell align="left">{row.retrievalItemInstruction && row.retrievalItemInstruction.retrievalPlatform ? row.retrievalItemInstruction.retrievalPlatform  :''}</TableCell>
                <TableCell align="left">
                    <Button
                        fullWidth
                        onClick={e => handleEditRetrieval(e, row)}
                        variant="contained"
                    >
                        {'Edit'}
                    </Button>
                </TableCell>
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
                                        <TableCell component="th" scope="row" style={{ width: '50%', whiteSpace: 'pre-wrap' }}>
                                                {
                                                    row.retrievalItemInstruction && row.retrievalItemInstruction.instructions && prettyPrintJSONFromBytes(row.retrievalItemInstruction.instructions)
                                                }
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

export const prettyPrintJSON = (json: any): string => {
    try {
        // Check if the input is a string that needs to be parsed
        if (typeof json === 'string') {
            json = JSON.parse(json);
        }
        return JSON.stringify(json, null, 2)

    } catch (error) {
        console.error('Error parsing or formatting JSON:', error);
        return ''; // Return an empty string in case of error
    }
};

export const prettyPrintJSONFromBytes = (byteArrayString: string | undefined): string => {
    if (byteArrayString === undefined || byteArrayString === '') {
        return '';
    }
    try {
        // Assuming byteArrayString is a base64 encoded string of the byte array
        const decodedString = atob(byteArrayString);
        let jsonObject = JSON.parse(decodedString);
        delete(jsonObject['instructions'])
        return JSON.stringify(jsonObject, null, 2); // Pretty print with 2 spaces indentation
    } catch (error) {
        console.error('Error parsing JSON:', error);
        return byteArrayString; // Fallback to original string in case of error
    }
};

import * as React from "react";
import {TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import {Assistant} from "../../redux/ai/ai.types2";


export function AssistantsRow(props: { row: Assistant, index: number, handleClick: any, checked: boolean}) {
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
                <TableCell align="left">{row.id}</TableCell>
                <TableCell align="left">{row.name? row.name : 'Null'}</TableCell>
                <TableCell align="left">{row.model}</TableCell>
            </TableRow>
            {/*<TableRow>*/}
            {/*    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>*/}
            {/*        <Collapse in={open} timeout="auto" unmountOnExit>*/}
            {/*            <Box sx={{ margin: 1 }}>*/}
            {/*                <Typography variant="h6" gutterBottom component="div">*/}
            {/*                    Retrieval Details*/}
            {/*                </Typography>*/}
            {/*                <Table size="small" aria-label="purchases">*/}
            {/*                    <TableBody>*/}
            {/*                        <TableRow >*/}
            {/*                            <TableCell component="th" scope="row" style={{ width: '50%', whiteSpace: 'pre-wrap' }}>*/}
            {/*                            </TableCell>*/}
            {/*                        </TableRow>*/}
            {/*                    </TableBody>*/}
            {/*                </Table>*/}
            {/*            </Box>*/}
            {/*        </Collapse>*/}
            {/*    </TableCell>*/}
            {/*</TableRow>*/}
        </React.Fragment>
    );
}

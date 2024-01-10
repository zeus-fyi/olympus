import * as React from "react";
import {Collapse, TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import {EvalFn} from "../../redux/ai/ai.types";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";

export function EvalRow(props: { row: EvalFn, index: number, handleClick: any, checked: boolean}) {
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
                <TableCell component="th" scope="row">
                    {row.evalID? row.evalID : 0}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.evalName}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.evalGroupName}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.evalType}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.evalFormat}
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={11}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Eval Metrics Details
                            </Typography>
                            <Table size="small" aria-label="sub-analysis">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Metric Name</TableCell>
                                        {/*<TableCell>Operator</TableCell>*/}
                                        {/*<TableCell>Eval State</TableCell>*/}
                                        <TableCell>Description</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.evalMetrics && row.evalMetrics.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell>{data.evalMetricName}</TableCell>
                                            {/*<TableCell>{data.evalOperator}</TableCell>*/}
                                            {/*<TableCell>{data.evalState}</TableCell>*/}
                                            <TableCell>{data.evalModelPrompt}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}

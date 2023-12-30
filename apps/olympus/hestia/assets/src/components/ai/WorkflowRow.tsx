import * as React from "react";
import {Box, Collapse, TableRow, Typography} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import {WorkflowTemplate} from "../../redux/ai/ai.types";
import TableHead from "@mui/material/TableHead";

export function WorkflowRow(props: { row: WorkflowTemplate, index: number, handleClick: any, checked: boolean}) {
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
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={11}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Workflow Details
                            </Typography>
                            <Table size="small" aria-label="sub-analysis">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Task Name</TableCell>
                                        <TableCell>Task Type</TableCell>
                                        <TableCell>Cycle Count</TableCell>
                                        <TableCell style={{ width: '15%'}}>Model</TableCell>
                                        <TableCell style={{ width: '50%', whiteSpace: 'pre-wrap' }}>Prompt</TableCell>
                                        <TableCell>Retrieval Name</TableCell>
                                        <TableCell>Retrieval Platform</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.tasks && row.tasks.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell>{data.taskName}</TableCell>
                                            <TableCell>{data.taskType}</TableCell>
                                            <TableCell>{data.cycleCount}</TableCell>
                                            <TableCell style={{ width: '15%'}}>{data.model}</TableCell>
                                            <TableCell style={{ width: '50%', whiteSpace: 'pre-wrap' }}>
                                                {data.prompt}
                                            </TableCell>
                                            <TableCell>{data.retrievalName ? data.retrievalName : 'analysis-aggregation'}</TableCell>
                                            <TableCell>{data.retrievalPlatform ? data.retrievalPlatform : 'analysis-aggregation'}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                            <Box sx={{ margin: 1 }}>
                                <Typography variant="h6" gutterBottom component="div">
                                    Eval Details
                                </Typography>
                            </Box>
                            <Table>
                                <TableBody>
                                    {row.tasks && row.tasks.map((task, taskIndex) => (
                                        task.evalFns && task.evalFns.map((evalFn, evalFnIndex) => (
                                            <TableRow key={evalFnIndex}>
                                                <TableCell>{evalFn.evalName}</TableCell>
                                                <TableCell>{evalFn.evalGroupName}</TableCell>
                                                <TableCell>{evalFn.evalType}</TableCell>
                                                <TableCell>{evalFn.evalModel}</TableCell>
                                                <TableCell>{evalFn.cycleCount}</TableCell>
                                                <TableCell>{evalFn.evalFormat}</TableCell>
                                            </TableRow>
                                        ))
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


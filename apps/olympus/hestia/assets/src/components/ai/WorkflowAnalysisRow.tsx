import * as React from "react";
import {Box, Collapse, TableRow, Typography} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import {OrchestrationsAnalysis} from "../../redux/ai/ai.types";
import TableHead from "@mui/material/TableHead";
import {prettyPrintJSON} from "./RetrievalsRow";

export function WorkflowAnalysisRow(props: { row: OrchestrationsAnalysis, index: number, handleClick: any, checked: boolean}) {
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
                <TableCell align="left">{row.orchestration.orchestrationID}</TableCell>
                <TableCell align="left">{row.orchestration.orchestrationName}</TableCell>
                <TableCell align="left">{row.orchestration.groupName}</TableCell>
                <TableCell align="left">{row.orchestration.type}</TableCell>
                <TableCell align="left">{row.orchestration.active ? 'Yes' : 'No'}</TableCell>
                <TableCell align="left">{row.runCycles}</TableCell>
                <TableCell align="left">{row.totalWorkflowTokenUsage}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={12}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Run Details
                            </Typography>
                            <Table size="small" aria-label="sub-analysis">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Task ID</TableCell>
                                        <TableCell>Task Name</TableCell>
                                        <TableCell>Task Type</TableCell>
                                        <TableCell>Cycle</TableCell>
                                        <TableCell>Start</TableCell>
                                        <TableCell>End</TableCell>
                                        <TableCell style={{ width: '15%'}}>Model</TableCell>
                                        <TableCell>Prompt Tokens</TableCell>
                                        <TableCell>Completion Tokens</TableCell>
                                        <TableCell>Total Tokens</TableCell>
                                        <TableCell style={{ width: '10%', whiteSpace: 'pre-wrap' }}>Prompt</TableCell>
                                        <TableCell style={{ width: '20%', whiteSpace: 'pre-wrap' }}>Completion Choices</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.aggregatedData && row.aggregatedData.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell>{data.sourceTaskId}</TableCell>
                                            <TableCell>{data.taskName}</TableCell>
                                            <TableCell>{data.taskType}</TableCell>
                                            <TableCell>{data.runningCycleNumber}</TableCell>
                                            <TableCell>{data.searchWindowUnixStart}</TableCell>
                                            <TableCell>{data.searchWindowUnixEnd}</TableCell>
                                            <TableCell style={{ width: '15%'}}>{data.model}</TableCell>
                                            <TableCell>{data.promptTokens}</TableCell>
                                            <TableCell>{data.completionTokens}</TableCell>
                                            <TableCell>{data.totalTokens}</TableCell>
                                            <TableCell style={{ width: '15%', whiteSpace: 'pre-wrap' }}>
                                                {data.prompt !== undefined ? prettyPrintJSON(data.prompt) : ""}
                                            </TableCell>
                                            <TableCell style={{ width: '15%', whiteSpace: 'pre-wrap' }}>
                                                {data.completionChoices !== undefined ? prettyPrintJSON(data.completionChoices) : ""}
                                            </TableCell>
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


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
            {row.aggregatedEvalResults && row.aggregatedEvalResults.length > 0 && (
                <TableRow>
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={12}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box sx={{ margin: 1 }}>
                                <Typography variant="h6" gutterBottom component="div">
                                    Eval Results
                                </Typography>
                                <Table size="small" aria-label="eval-results">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Eval Name</TableCell>
                                            <TableCell>Metric Name</TableCell>
                                            <TableCell>Metric ID</TableCell>
                                            <TableCell>Result ID</TableCell>
                                            <TableCell>Metric Result</TableCell>
                                            <TableCell>Comparison Boolean</TableCell>
                                            <TableCell>Comparison Number</TableCell>
                                            <TableCell>Comparison String</TableCell>
                                            <TableCell>Metric Data Type</TableCell>
                                            <TableCell>Operator</TableCell>
                                            <TableCell>State</TableCell>
                                            <TableCell>Running Cycle Number</TableCell>
                                            <TableCell>Start Unix Time</TableCell>
                                            <TableCell>End Unix Time</TableCell>
                                            <TableCell>Result Outcome</TableCell>
                                            <TableCell>Metadata</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {row.aggregatedEvalResults.map((evalResult, evalIndex) => (
                                            <TableRow key={evalIndex}>
                                                <TableCell>{evalResult.evalName}</TableCell>
                                                <TableCell>{evalResult.evalMetricName}</TableCell>
                                                <TableCell>{evalResult.evalMetricID}</TableCell>
                                                <TableCell>{evalResult.evalMetricsResultId}</TableCell>
                                                <TableCell>{evalResult.evalMetricResult}</TableCell>
                                                <TableCell>{evalResult.evalComparisonBoolean ? 'True' : 'False'}</TableCell>
                                                <TableCell>{evalResult.evalComparisonNumber}</TableCell>
                                                <TableCell>{evalResult.evalComparisonString}</TableCell>
                                                <TableCell>{evalResult.evalMetricDataType}</TableCell>
                                                <TableCell>{evalResult.evalOperator}</TableCell>
                                                <TableCell>{evalResult.evalState}</TableCell>
                                                <TableCell>{evalResult.runningCycleNumber}</TableCell>
                                                <TableCell>{evalResult.searchWindowUnixStart}</TableCell>
                                                <TableCell>{evalResult.searchWindowUnixEnd}</TableCell>
                                                <TableCell>{evalResult.evalResultOutcome ? 'Pass' : 'Fail'}</TableCell>
                                                <TableCell>{evalResult.evalMetadata}</TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </Box>
                        </Collapse>
                    </TableCell>
                </TableRow>
            )}
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


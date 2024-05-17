import * as React from "react";
import {Box, Collapse, LinearProgress, TableRow, Typography} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import {OrchestrationsAnalysis} from "../../redux/ai/ai.types.runs";
import TableHead from "@mui/material/TableHead";
import {prettyPrintWfRunRowJSON} from "../ai/RetrievalsRow";
import CsvExportButton from "./CsvHelpers";

function convertUnixNanoToLocalTimeString(unixNano: string): string {
    // Convert the string to a BigInt to handle large numbers, then to a number in milliseconds
    const milliseconds = Number(BigInt(unixNano) / BigInt(1e6));

    // Create a Date object
    const date = new Date(milliseconds);

    // Convert to local time string
    return date.toLocaleString();
}

export function WorkflowAnalysisRow(props: { row: OrchestrationsAnalysis, index: number, handleClick: any, checked: boolean, csvExport: boolean; open : boolean; handleOpen: any, isAdmin: boolean }) {
    const {csvExport, row, index, handleClick, checked, open, handleOpen, isAdmin } = props;

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => handleOpen(index)}
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
                <TableCell align="left">{row.orchestration.orchestrationStrID}</TableCell>
                <TableCell align="left">{convertUnixNanoToLocalTimeString(row.orchestration.orchestrationStrID)}</TableCell>
                <TableCell style={{ fontWeight: 'normal', color: 'white', minWidth: 50}}>
                    <LinearProgress variant="determinate" value={row.progress} />
                </TableCell>
                <TableCell align="left">{row.orchestration.orchestrationName}</TableCell>
                <TableCell align="left">{row.orchestration.groupName}</TableCell>
                <TableCell align="left">{row.orchestration.type}</TableCell>
                <TableCell align="left">{row.orchestration.active ? 'Yes' : 'No'}</TableCell>
                <TableCell align="left">{row.runCycles}</TableCell>
                <TableCell align="left">{row.totalWorkflowTokenUsage}</TableCell>
                { csvExport && <CsvExportButton name={row.orchestration.groupName} orchStrID={row.orchestration.orchestrationStrID} results={row} isAdmin={isAdmin}/>}
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
                                            <TableCell>Eval Metric Result ID</TableCell>
                                            <TableCell>Eval Name</TableCell>
                                            <TableCell>Cycle</TableCell>
                                            <TableCell>Iteration</TableCell>
                                            <TableCell>Field Name</TableCell>
                                            <TableCell>Data Type</TableCell>
                                            <TableCell>State</TableCell>
                                            <TableCell>Operator</TableCell>
                                            <TableCell>Expected</TableCell>
                                            <TableCell>Actual</TableCell>
                                            <TableCell>Details</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {row.aggregatedEvalResults && row.aggregatedEvalResults.map((evalResult, evalIndex) => {
                                            return (
                                                <TableRow key={evalIndex}>
                                                    <TableCell>{evalResult.evalMetricResult ? evalResult.evalMetricResult.evalMetricsResultStrID : ''}</TableCell>
                                                    <TableCell>{evalResult.evalName}</TableCell>
                                                    <TableCell>
                                                        {evalResult.evalMetricResult && evalResult.evalMetricResult.runningCycleNumber
                                                            ? evalResult.evalMetricResult.runningCycleNumber
                                                            : 0}
                                                    </TableCell>
                                                    <TableCell>
                                                        {evalResult.evalMetricResult && evalResult.evalMetricResult.evalIterationCount
                                                            ? evalResult.evalMetricResult.evalIterationCount
                                                            : 0}
                                                    </TableCell>
                                                    <TableCell>
                                                        {evalResult.evalField && evalResult.evalField.fieldName
                                                            ?  evalResult.evalField.fieldName
                                                            : 'No name available'}
                                                    </TableCell>
                                                    <TableCell>
                                                        {evalResult.evalField && evalResult.evalField.dataType
                                                            ?  evalResult.evalField.dataType
                                                            : 'No type available'}
                                                    </TableCell>
                                                    <TableCell>{evalResult.evalState}</TableCell>
                                                    <TableCell>{evalResult.evalOperator}</TableCell>
                                                    <TableCell>
                                                        {evalResult.evalExpectedResultState}
                                                    </TableCell>
                                                    <TableCell>
                                                        {evalResult.evalMetricResult && evalResult.evalMetricResult.evalResultOutcomeStateStr
                                                            ? evalResult.evalMetricResult.evalResultOutcomeStateStr
                                                            : 'No state available'}
                                                    </TableCell>
                                                    <TableCell>
                                                        {evalResult.evalMetricResult && evalResult.evalMetricResult.evalMetadata
                                                            ? evalResult.evalMetricResult.evalMetadata.evalOpCtxStr
                                                            : ''}
                                                    </TableCell>
                                                </TableRow>
                                            );
                                        })}
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
                                Retrieval Details
                            </Typography>
                            <Table size="small" aria-label="sub-analysis">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Result ID</TableCell>
                                        <TableCell>Retrieval Name</TableCell>
                                        <TableCell>Cycle</TableCell>
                                        <TableCell>Iteration</TableCell>
                                        <TableCell>Offset</TableCell>
                                        <TableCell>Status</TableCell>
                                        <TableCell>Usage</TableCell>
                                        {/*<TableCell>Start</TableCell>*/}
                                        <TableCell>End</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.aggregatedRetrievalResults && row.aggregatedRetrievalResults.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell>{data.workflowResultStrID}</TableCell>
                                            <TableCell>{data.retrievalName}</TableCell>
                                            <TableCell>{data.runningCycleNumber}</TableCell>
                                            <TableCell>{data.iterationCount}</TableCell>
                                            <TableCell>{data.chunkOffset}</TableCell>
                                            <TableCell>{data.status}</TableCell>
                                            <TableCell>{data.skipRetrieval ? 'skipped' : 'used'}</TableCell>
                                            {/*<TableCell>{data.searchWindowUnixStart}</TableCell>*/}
                                            <TableCell>{convertUnixNanoToLocalTimeString(data.workflowResultStrID)}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
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
                                        <TableCell>Result ID</TableCell>
                                        <TableCell>Task ID</TableCell>
                                        <TableCell>Task Name</TableCell>
                                        <TableCell>Task Type</TableCell>
                                        <TableCell>Cycle</TableCell>
                                        <TableCell>Iteration</TableCell>
                                        <TableCell>Offset</TableCell>
                                        <TableCell>Usage</TableCell>
                                        <TableCell>Start</TableCell>
                                        <TableCell>End</TableCell>
                                        <TableCell>Model</TableCell>
                                        <TableCell>Prompt</TableCell>
                                        <TableCell>Completion Tokens</TableCell>
                                        <TableCell>Total Tokens</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.aggregatedData && row.aggregatedData.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell>{data.workflowResultID}</TableCell>
                                            <TableCell>{data.sourceTaskID}</TableCell>
                                            <TableCell>{data.taskName}</TableCell>
                                            <TableCell>{data.taskType}</TableCell>
                                            <TableCell>{data.runningCycleNumber}</TableCell>
                                            <TableCell>{data.iterationCount}</TableCell>
                                            <TableCell>{data.chunkOffset}</TableCell>
                                            <TableCell>{data.skipAnalysis ? 'skipped' : 'used'}</TableCell>
                                            <TableCell>{data.searchWindowUnixStart}</TableCell>
                                            <TableCell>{data.searchWindowUnixEnd}</TableCell>
                                            <TableCell>{data.model}</TableCell>
                                            <TableCell>{data.promptTokens}</TableCell>
                                            <TableCell>{data.completionTokens}</TableCell>
                                            <TableCell>{data.totalTokens}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                            <Table  sx={{ mb: 4, mt: 4}}  aria-label="sub-analysis">
                                <TableRow>
                                    <TableCell >Result ID</TableCell>
                                    <TableCell >Prompt</TableCell>
                                    <TableCell >Completion Choices</TableCell>
                                </TableRow>
                                <TableBody>
                                    {row.aggregatedData && row.aggregatedData.map((data, dataIndex) => (
                                        <TableRow key={dataIndex}>
                                            <TableCell style={{ maxWidth: 120  }} >
                                                {data.workflowResultID}
                                            </TableCell>
                                            <TableCell style={{ whiteSpace: 'pre-wrap', wordWrap: 'break-word',  maxWidth: 400  }}>
                                                {data && data.prompt !== undefined ? prettyPrintWfRunRowJSON(data.prompt): ""}
                                            </TableCell>
                                            <TableCell style={{ whiteSpace: 'pre-wrap', wordWrap: 'break-word', maxWidth: 400  }}>
                                                {data && data.completionChoices !== undefined ? prettyPrintWfRunRowJSON(data.completionChoices) : ""}
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


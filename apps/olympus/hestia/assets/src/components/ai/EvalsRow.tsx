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
import Button from "@mui/material/Button";
import {useDispatch} from "react-redux";
import {setEval} from "../../redux/ai/ai.reducer";

export function EvalRow(props: { row: EvalFn, index: number, handleClick: any, checked: boolean}) {
    const { row, index, handleClick, checked } = props;
    const [open, setOpen] = React.useState(false);
    const dispatch = useDispatch();
    const handleEditEvalFunction = async (e: any, ef: EvalFn) => {
        e.preventDefault();
        dispatch(setEval(ef))
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
                <TableCell align="left">
                    <Button
                        fullWidth
                        onClick={e => handleEditEvalFunction(e, row)}
                        variant="contained"
                    >
                        {'Edit'}
                    </Button>
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
                                        <TableCell>Metric ID</TableCell>
                                        <TableCell>Metric Name</TableCell>
                                        <TableCell>Description</TableCell>
                                        <TableCell>Operator</TableCell>
                                        <TableCell>Eval State</TableCell>
                                        <TableCell>Expected Result</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.evalMetrics && row.evalMetrics.map((data, dataIndex) => (
                                        <TableRow key={data.evalMetricName}>
                                            <TableCell>{data.evalMetricID ? data.evalMetricID : 'N/A'}</TableCell>
                                            <TableCell>{data.evalMetricName}</TableCell>
                                            <TableCell>{data.evalModelPrompt}</TableCell>
                                            <TableCell>{data.evalOperator}</TableCell>
                                            <TableCell>{data.evalState}</TableCell>
                                            <TableCell>{data.evalMetricResult}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        {row.triggerFunctions && row.triggerFunctions.length > 0 && (
                            <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Eval Triggers Details
                            </Typography>
                            <Table size="small" aria-label="sub-analysis">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Trigger Name</TableCell>
                                        <TableCell>Trigger Group</TableCell>
                                        <TableCell>Eval State</TableCell>
                                        <TableCell>Trigger On</TableCell>
                                        <TableCell>Output Env</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.triggerFunctions && row.triggerFunctions.map((data, dataIndex) => (
                                        data.evalTriggerActions && data.evalTriggerActions.map((evalTrigger, triggerIndex) => (
                                            <TableRow key={'0' + '-' + dataIndex + '-' + triggerIndex}>
                                                <TableCell>{data.triggerName}</TableCell>
                                                <TableCell>{data.triggerGroup}</TableCell>
                                                <TableCell>{evalTrigger.evalTriggerState}</TableCell>
                                                <TableCell>{evalTrigger.evalResultsTriggerOn}</TableCell>
                                                <TableCell>{data.triggerEnv}</TableCell>
                                            </TableRow>
                                        ))
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                            )}
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}

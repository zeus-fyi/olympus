import * as React from "react";
import {Collapse, TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import {TriggerActionsApproval} from "../../redux/ai/ai.types.retrievals";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import Button from "@mui/material/Button";

// handleClick: any, checked: boolean
export function ActionApprovalsRow(props: any) {
    const { row, handleActionApprovalRequest, index} = props;
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
                {/*<TableCell align="center" >*/}
                {/*    <Checkbox*/}
                {/*        checked={checked}*/}
                {/*        onChange={() => handleClick(index)}*/}
                {/*        color="primary"*/}
                {/*    />*/}
                {/*</TableCell>*/}
                {/*<TableCell component="th" scope="row">*/}
                {/*    {row.triggerID ? row.triggerID : 0}*/}
                {/*</TableCell>*/}
                <TableCell component="th" scope="row">
                    {row.triggerGroup}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.triggerName}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.triggerAction}
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={11}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        {row.triggerActionsApprovals && row.triggerActionsApprovals.length > 0 && (
                            <Box sx={{ margin: 1 }}>
                                <Typography variant="h6" gutterBottom component="div">
                                    Action Details
                                </Typography>
                                <Table size="small" aria-label="sub-analysis">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Workflow Result ID</TableCell>
                                            <TableCell>Approval ID</TableCell>
                                            <TableCell>Approval State</TableCell>
                                            <TableCell>Request Summary</TableCell>
                                            <TableCell>Updated At</TableCell>
                                            <TableCell></TableCell>
                                            <TableCell></TableCell>
                                            <TableCell></TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {row.triggerActionsApprovals && row.triggerActionsApprovals
                                            .filter((data: TriggerActionsApproval) => data.approvalState === 'pending')
                                            .map((data: TriggerActionsApproval, dataIndex: number) => (
                                                <TableRow key={dataIndex}>
                                                    <TableCell>{data.workflowResultStrID}</TableCell>
                                                    <TableCell>{data.approvalStrID}</TableCell>
                                                    <TableCell>{capitalizeFirstLetter(data.approvalState)}</TableCell>
                                                    <TableCell>{data.requestSummary}</TableCell>
                                                    <TableCell>{new Date(data.updatedAt).toLocaleString()}</TableCell>
                                                    <TableCell align="left"><Button onClick={event => handleActionApprovalRequest(event,'approved', data)} fullWidth variant="contained" >{'Approve'}</Button></TableCell>
                                                    <TableCell align="left"><Button onClick={event => handleActionApprovalRequest(event,'rejected', data)} fullWidth variant="contained" >{'Reject'}</Button></TableCell>
                                                </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </Box>
                        )}
                    </Collapse>
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={11}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        {row.triggerActionsApprovals && row.triggerActionsApprovals.length > 0 && (
                            <Box sx={{ margin: 1 }}>
                                <Typography variant="h6" gutterBottom component="div">
                                    Historical Actions
                                </Typography>
                                <Table size="small" aria-label="sub-analysis">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Workflow Result ID</TableCell>
                                            <TableCell>Approval ID</TableCell>
                                            <TableCell>Final State</TableCell>
                                            <TableCell>Request Summary</TableCell>
                                            <TableCell>Updated At</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {row.triggerActionsApprovals && row.triggerActionsApprovals
                                            .filter((data: TriggerActionsApproval) => data.approvalState !== 'pending')
                                            .map((data: TriggerActionsApproval, dataIndex: number) => (
                                                <TableRow key={dataIndex}>
                                                    <TableCell>{data.workflowResultStrID}</TableCell>
                                                    <TableCell>{data.approvalStrID}</TableCell>
                                                    <TableCell>{capitalizeFirstLetter(data.approvalState)}</TableCell>
                                                    <TableCell>{data.requestSummary}</TableCell>
                                                    <TableCell>{new Date(data.updatedAt).toLocaleString()}</TableCell>
                                                </TableRow>
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
function capitalizeFirstLetter(v: string) {
    if (v === null || v === undefined || v.length === 0) {
        return v;
    }
    return v.charAt(0).toUpperCase() + v.slice(1);
}

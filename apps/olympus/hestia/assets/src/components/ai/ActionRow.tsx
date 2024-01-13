import * as React from "react";
import {Collapse, TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Checkbox from "@mui/material/Checkbox";
import {TriggerAction} from "../../redux/ai/ai.types2";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import Button from "@mui/material/Button";
import {setTriggerAction} from "../../redux/ai/ai.reducer";
import {useDispatch} from "react-redux";

export function ActionRow(props: { row: TriggerAction, index: number, handleClick: any, checked: boolean}) {
    const { row, index, handleClick, checked } = props;
    const [open, setOpen] = React.useState(false);
    const dispatch = useDispatch();

    const handleEditTriggerAction = async (e: any, ta: TriggerAction) => {
        e.preventDefault();
        dispatch(setTriggerAction(ta))
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
                    {row.triggerID ? row.triggerID : 0}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.triggerGroup}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.triggerName}
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.triggerEnv}
                </TableCell>
                <TableCell align="left">
                    <Button onClick={event => handleEditTriggerAction(event, row)} fullWidth variant="contained" >{'Edit'}</Button>
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={11}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        {row.evalTriggerActions && row.evalTriggerActions.length > 0 && (
                            <Box sx={{ margin: 1 }}>
                                <Typography variant="h6" gutterBottom component="div">
                                    Eval Triggers Details
                                </Typography>
                                <Table size="small" aria-label="sub-analysis">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Eval State</TableCell>
                                            <TableCell>Trigger On</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {row.evalTriggerActions && row.evalTriggerActions.map((data, dataIndex) => (
                                                <TableRow key={dataIndex}>
                                                    <TableCell>{data.evalTriggerState}</TableCell>
                                                    <TableCell>{data.evalResultsTriggerOn}</TableCell>
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


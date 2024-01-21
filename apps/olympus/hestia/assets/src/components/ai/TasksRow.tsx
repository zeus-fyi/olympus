import * as React from "react";
import {Box, Collapse, TableRow, Typography} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import Checkbox from "@mui/material/Checkbox";
import Button from "@mui/material/Button";
import {TaskModelInstructions} from "../../redux/ai/ai.types";
import {setEditAggregateTask, setEditAnalysisTask} from "../../redux/ai/ai.reducer";
import {useDispatch} from "react-redux";
import TableHead from "@mui/material/TableHead";

export function TasksRow(props: { row: TaskModelInstructions, index: number, handleClick: any,checked: boolean}) {
    const { row, index, handleClick, checked } = props;
    const [open, setOpen] = React.useState(false);
    const dispatch = useDispatch();
    const handleEditClick = () => {
        if (row.taskType === 'analysis') {
            dispatch(setEditAnalysisTask(row));
        } else if (row.taskType === 'aggregation') {
            dispatch(setEditAggregateTask(row));
        }
    };
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
                <TableCell align="left">{row.taskID}</TableCell>
                <TableCell align="left">{row.taskGroup}</TableCell>
                <TableCell align="left">{row.taskName}</TableCell>
                <TableCell align="left">{row.model}</TableCell>
                <TableCell align="left">{row.responseFormat}</TableCell>
                <TableCell align="left"><Button fullWidth variant="contained" onClick={handleEditClick}>Edit</Button></TableCell>
            </TableRow>
            <TableRow>
                {
                    row.responseFormat === 'json' || row.responseFormat === 'social-media-engagement' ?
                        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>
                            <Collapse in={open} timeout="auto" unmountOnExit>
                                {row.schemas && row.schemas.length > 0 && (
                                    <Box sx={{ margin: 1 }}>
                                        <Typography variant="h6" gutterBottom component="div">
                                            Schemas
                                        </Typography>
                                        <Table size="small" aria-label="sub-analysis">
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell>Group</TableCell>
                                                    <TableCell>Name</TableCell>
                                                    <TableCell>Fields</TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {row.schemas && row.schemas.map((data, dataIndex) => (
                                                    data &&
                                                    <TableRow key={dataIndex}>
                                                        <TableCell>{data.schemaGroup}</TableCell>
                                                        <TableCell>{data.schemaName}</TableCell>
                                                        {data.fields && data.fields.length > 0 && data.fields.map((field, fieldIndex) => (
                                                            <TableRow key={fieldIndex}>
                                                                <TableCell>{field.dataType}</TableCell>
                                                                <TableCell>{field.fieldName}</TableCell>
                                                                {/*<TableCell>{field.fieldDescription}</TableCell>*/}
                                                            </TableRow>
                                                        ))}
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </Box>
                                )}
                            </Collapse>
                        </TableCell>
                        :
                        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>
                            <Collapse in={open} timeout="auto" unmountOnExit>
                                <Box sx={{ margin: 1 }}>
                                    <Typography variant="h6" gutterBottom component="div">
                                        Prompt
                                    </Typography>
                                    <Table size="small" aria-label="prompt">
                                        <TableBody>
                                            <TableRow >
                                                <TableCell component="th" scope="row">
                                                    {row.prompt}
                                                </TableCell>
                                            </TableRow>
                                        </TableBody>
                                    </Table>
                                </Box>
                            </Collapse>
                        </TableCell>
                }
            </TableRow>
        </React.Fragment>
    );
}

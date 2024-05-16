import * as React from "react";
import {TableRow} from "@mui/material";
import TableCell from "@mui/material/TableCell";
import {UserFlowStats} from "../../../../redux/flows/flows.actions";


function convertUnixNanoToLocalTimeString(unixNano: string): string {
    // Convert the string to a BigInt to handle large numbers, then to a number in milliseconds
    const milliseconds = Number(BigInt(unixNano) / BigInt(1e6));
    // Create a Date object
    const date = new Date(milliseconds);
    // Convert to local time string
    return date.toLocaleString();
}

// handleClick: any, checked: boolean,
export function UserStatsRow(props: { row: UserFlowStats, index: number, csvExport: boolean; open : boolean; handleOpen: any }) {
    const {csvExport, row, index, open, handleOpen } = props;

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                {/*<TableCell>*/}
                {/*    <IconButton*/}
                {/*        aria-label="expand row"*/}
                {/*        size="small"*/}
                {/*        onClick={() => handleOpen(index)}*/}
                {/*    >*/}
                {/*        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}*/}
                {/*    </IconButton>*/}
                {/*</TableCell>*/}
                {/*<TableCell align="center" >*/}
                {/*    <Checkbox*/}
                {/*        checked={checked}*/}
                {/*        onChange={() => handleClick(index)}*/}
                {/*        color="primary"*/}
                {/*    />*/}
                {/*</TableCell>*/}
                {/*<TableCell align="left">{convertUnixNanoToLocalTimeString(row.orchestration.orchestrationStrID)}</TableCell>*/}
                {/*<TableCell style={{ fontWeight: 'normal', color: 'white', minWidth: 50}}>*/}
                {/*    <LinearProgress variant="determinate" value={row.progress} />*/}
                {/*</TableCell>*/}
                <TableCell align="left">{row.orgID}</TableCell>
                <TableCell align="left">{convertUnixNanoToLocalTimeString(row.orgID)}</TableCell>
                <TableCell align="left">{row.email}</TableCell>
                <TableCell align="left">{row.flowCount}</TableCell>

            </TableRow>
        </React.Fragment>
    );
}


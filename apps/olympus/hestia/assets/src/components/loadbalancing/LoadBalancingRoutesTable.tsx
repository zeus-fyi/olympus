import * as React from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Checkbox from "@mui/material/Checkbox";
import {TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

export function LoadBalancingRoutesTable(props: any) {
    const { loading,selected, endpoints, handleSelectAllClick, handleClick } = props

    if (loading) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }

    if (endpoints === null || endpoints === undefined) {
        return (<div></div>)
    }
    return (
        <div>
            <Box sx={{ mt: 4, mb: 4 }}>
                {selected.length > 0 && (
                    <Box sx={{ mb: 2 }}>
                        <span>({selected.length} selected endpoints)</span>
                        <Button variant="outlined" color="secondary" style={{marginLeft: '10px'}}>
                            Delete
                        </Button>
                    </Box>
                )}
                <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 650 }} aria-label="simple table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333'}} >
                                <TableCell padding="checkbox">
                                    <Checkbox
                                        color="primary"
                                        indeterminate={selected.length > 0 && selected.length < endpoints.length}
                                        checked={endpoints.length > 0 && selected.length === endpoints.length}
                                        onChange={handleSelectAllClick}
                                    />
                                </TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Endpoints</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {endpoints.map((row: string, i: number) => (
                                <TableRow
                                    key={i}
                                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                >
                                    <TableCell padding="checkbox">
                                        <Checkbox
                                            checked={selected.indexOf(row) !== -1}
                                            onChange={() => handleClick(row)}
                                            color="primary"
                                        />
                                    </TableCell>
                                    <TableCell component="th" scope="row">
                                        {row}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        </div>
    );
}

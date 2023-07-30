import {useParams} from "react-router-dom";
import * as React from "react";
import {useEffect, useState} from "react";
import Box from "@mui/material/Box";
import {TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import {loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {useDispatch, useSelector} from "react-redux";
import {setEndpoints} from "../../redux/loadbalancing/loadbalancing.reducer";
import {RootState} from "../../redux/store";

export function LoadBalancingRoutesTable() {
    const params = useParams();
    const dispatch = useDispatch();
    const [loading, setLoading] = useState(true);
    const endpoints = useSelector((state: RootState) => state.loadBalancing.routes);

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await loadBalancingApiGateway.getEndpoints();
                const endpoints = response.data;
                dispatch(setEndpoints(endpoints.routes));
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
            }}
        fetchData(params);
    }, []);
    if (loading) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }

    if (endpoints === null || endpoints === undefined) {
        return (<div></div>)
    }

    return (
        <div>
            <Box sx={{ mt: 4, mb: 4 }}>
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                    <TableHead>
                        <TableRow style={{ backgroundColor: '#333'}} >
                            <TableCell style={{ color: 'white'}} align="left">Endpoints</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {endpoints.map((row: any, i: number) => (
                            <TableRow
                                key={i}
                                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                            >
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
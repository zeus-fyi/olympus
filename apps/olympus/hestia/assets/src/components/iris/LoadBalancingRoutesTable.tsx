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

function LoadBalancingRoutesTable(cluster: any) {
    const params = useParams();
    const [activeClusterTopologies, setActiveClusterTopologies] = useState([{}]);
    const [statusMessage, setStatusMessage] = useState('');
    const [statusMessageRowIndex, setStatusMessageRowIndex] = useState<number | null>(null);

    // const onClickRolloutUpgrade = async (index: number, clusterClassName: string) => {
    //     try {
    //         const response = await loadBalancingApiGateway.getEndpoints();
    //         const statusCode = response.status;
    //         if (statusCode === 202) {
    //             setStatusMessageRowIndex(index);
    //             setStatusMessage(`Cluster ${clusterClassName} update in progress`);
    //         } else if (statusCode === 200){
    //             setStatusMessageRowIndex(index);
    //             setStatusMessage(`Cluster ${clusterClassName} already up to date`);
    //         } else {
    //             setStatusMessageRowIndex(index);
    //             setStatusMessage(`Cluster ${clusterClassName} had an unexpected response: status code ${statusCode}`);
    //         }
    //     } catch (e) {
    //         setStatusMessageRowIndex(index);
    //         setStatusMessage(`Cluster ${clusterClassName} failed to update`);
    //     }
    // }

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await loadBalancingApiGateway.getEndpoints();
                const endpoints: any[] = response.data;
                // const clusterTopologyRows = clustersTopologyData.map((topology: any) =>
                //     createTopologyData(topology.topologyID, topology.clusterName, topology.componentBaseName, topology.skeletonBaseName),
                // );
                // setActiveClusterTopologies(clusterTopologyRows);
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(params);
    }, []);

    return (
        <div>
            <Box sx={{ mt: 4, mb: 4 }}>
                <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 650 }} aria-label="simple table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333'}} >
                                <TableCell style={{ color: 'white'}} align="left">Endpoint</TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                                <TableCell style={{ color: 'white'}} align="left"></TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {activeClusterTopologies
                                .filter(
                                    (item: any, index: number, self: any) =>
                                        index ===
                                        self.findIndex(
                                            (otherItem: any) => otherItem.clusterName === item.clusterName
                                        )
                                )
                                .map((row: any, i: number) => (
                                    <TableRow
                                        key={i}
                                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                    >
                                        <TableCell component="th" scope="row">
                                            {row.clusterName}
                                        </TableCell>
                                        {/*<TableCell align="left">*/}
                                        {/*    <Button onClick={() => onClickRolloutUpgrade(i, row.clusterName)} variant="contained">Deploy Latest</Button>*/}
                                        {/*    {statusMessageRowIndex === i && <div>{statusMessage}</div>}*/}
                                        {/*</TableCell>*/}
                                        <TableCell align="left">
                                            The Deploy Latest button will deploy the latest version configs to the cluster.
                                        </TableCell>
                                    </TableRow>
                                ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
            <TableContainer component={Paper}>
                {/*<Table sx={{ minWidth: 650 }} aria-label="simple table">*/}
                {/*    <TableHead>*/}
                {/*        <TableRow style={{ backgroundColor: '#333'}} >*/}
                {/*            <TableCell style={{ color: 'white'}}>TopologyID</TableCell>*/}
                {/*            <TableCell style={{ color: 'white'}} align="left">ClusterName</TableCell>*/}
                {/*            <TableCell style={{ color: 'white'}} align="left">ClusterBaseName</TableCell>*/}
                {/*            <TableCell style={{ color: 'white'}} align="left">SkeletonBaseName</TableCell>*/}
                {/*        </TableRow>*/}
                {/*    </TableHead>*/}
                {/*    <TableBody>*/}
                {/*        {activeClusterTopologies.map((row: any, i: number) => (*/}
                {/*            <TableRow*/}
                {/*                key={i}*/}
                {/*                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}*/}
                {/*            >*/}
                {/*                <TableCell component="th" scope="row">*/}
                {/*                    {row.topologyID}*/}
                {/*                </TableCell>*/}
                {/*                <TableCell align="left">{row.clusterName}</TableCell>*/}
                {/*                <TableCell align="left">{row.componentBaseName}</TableCell>*/}
                {/*                <TableCell align="left">{row.skeletonBaseName}</TableCell>*/}
                {/*            </TableRow>*/}
                {/*        ))}*/}
                {/*    </TableBody>*/}
                {/*</Table>*/}
            </TableContainer>
        </div>
    );
}
import * as React from "react";
import {useEffect, useState} from "react";
import {clustersApiGateway} from "../../gateway/clusters";
import {TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

function createTopologyData(
    topologyID: number,
    clusterName: string,
    clusterBaseName: string,
) {
    return {topologyID, clusterName, clusterBaseName};
}

function ActiveCloudClustersDeployedTopologiesContent(cluster: any) {
    const [activeClusterTopologies, setActiveClusterTopologies] = useState([{}]);
    useEffect(() => {
        const fetchData = async (cluster: any) => {
            try {
                const response = await clustersApiGateway.getClusterTopologies(cluster);
                const clustersTopologyData: any[] = response.data;
                const clusterTopologyRows = clustersTopologyData.map((topology: any) =>
                    createTopologyData(topology.topologyID, topology.clusterName, topology.clusterBaseName),
                );
                setActiveClusterTopologies(clusterTopologyRows);
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(cluster);
    }, []);
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>TopologyID</TableCell>
                        <TableCell align="left">ClusterName</TableCell>
                        <TableCell align="left">ClusterBaseName</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {activeClusterTopologies.map((row: any) => (
                        <TableRow
                            key={row.topologyID}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.topologyID}
                            </TableCell>
                            <TableCell align="left">{row.clusterName}</TableCell>
                            <TableCell align="left">{row.clusterBaseName}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

export default function ClustersPage() {
    return <ActiveCloudClustersDeployedTopologiesContent />;
}

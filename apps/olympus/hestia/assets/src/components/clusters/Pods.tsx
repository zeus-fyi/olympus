import {useParams} from "react-router-dom";
import * as React from "react";
import {useEffect, useState} from "react";
import {clustersApiGateway} from "../../gateway/clusters";
import {TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

function createPodsData(
    podName: string,
    podPhase: string
) {
    return {podName, podPhase};
}

export function PodsPageTable() {
    const params = useParams();
    const [pods, setPods] = useState([{}]);
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await clustersApiGateway.getClusterPodsAudit(params.id);
                const podSummaries = response.data.pods
                //console.log(podSummaries)
                let podsRows: any[] = [];
                for (const [key, value] of Object.entries(podSummaries)) {
                    let podInfo: any = value;
                    podsRows.push(createPodsData(key, podInfo.podPhase));
                }
                setPods(podsRows);
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(params);
    }, []);
    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ color: 'white'}}>PodName</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Status</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {pods.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.podName}
                            </TableCell>
                            <TableCell align="left">{row.podPhase}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

interface PodsSummaries {
    pods: {
        [key: string]: PodSummary;
    };
}

interface PodSummary {
    podName: string;
    phase: string;
    message: string;
    reason: string;
    startTime: string;
    podConditions: Array<{
        // fields for v1.PodCondition
        type: string;
        status: string;
        lastProbeTime: string;
        lastTransitionTime: string;
        reason: string;
        message: string;
    }>;
    initContainerStatuses: {
        [key: string]: {
            // fields for v1.ContainerStatus
            name: string;
            ready: boolean;
            restartCount: number;
            image: string;
            imageID: string;
            containerID: string;
            started: boolean;
        };
    };
    containerStatuses: {
        [key: string]: {
            // fields for v1.ContainerStatus
            name: string;
            ready: boolean;
            restartCount: number;
            image: string;
            imageID: string;
            containerID: string;
            started: boolean;
        };
    };
}
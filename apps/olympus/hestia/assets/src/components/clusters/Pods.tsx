import {useParams} from "react-router-dom";
import * as React from "react";
import {useEffect, useState} from "react";
import {clustersApiGateway} from "../../gateway/clusters";
import {Box, Button, FormControl, MenuItem, Select, TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import PodLogStreamClusterPage from "./PodLogStreamClusterPage";

function createPodsData(
    podName: string,
    podPhase: string,
    containers: string[],
) {
    return {podName, podPhase, containers};
}

export function PodsPageTable() {
    const params = useParams();
    const [pods, setPods] = useState([{}]);
    const [code, setCode] = useState('');
    const [selectedContainers, setSelectedContainers] = useState<Array<string>>([]);

    const handleContainerChange = (event: any, rowIndex: number) => {
        const selectedContainer = event.target.value as string;
        setSelectedContainers((prevSelectedContainers) => {
            const updatedContainers = [...prevSelectedContainers];
            updatedContainers[rowIndex] = selectedContainer;
            return updatedContainers;
        });
    };
    const onClickStreamLogs = async (podName: string, containerName: string) => {
        try {
            let res: any = await clustersApiGateway.getClusterPodLogs(params.id, podName, containerName)
            const statusCode = res.status;
            if (statusCode === 200 || statusCode === 204) {
                setCode(res.data)
            } else {
            }
        } catch (e) {
        }
    }
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await clustersApiGateway.getClusterPodsAudit(params.id);
                const podSummaries = response.data.pods
                let podsRows: any[] = [];
                for (const [key, value] of Object.entries(podSummaries)) {
                    let podInfo: any = value;
                    podsRows.push(createPodsData(key, podInfo.podPhase, podInfo.containers));
                }
                setPods(podsRows);
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(params).then(r =>
            console.log("")
        );
    }, []);
    useEffect(() => {
    }, [pods, selectedContainers]);

    return (
        <div>
            {pods && (
                <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 650 }} aria-label="simple table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333'}} >
                                <TableCell style={{ color: 'white'}}>PodName</TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Status</TableCell>
                                <TableCell style={{ color: 'white'}} align="left">Containers</TableCell>
                                <TableCell style={{ color: 'white'}} align="right"></TableCell>
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
                                    <TableCell align="left">
                                        {row.containers && (
                                            <FormControl fullWidth>
                                                <Select
                                                    labelId={`container-select-label-${i}`}
                                                    value={selectedContainers[i] || row.containers[row.containers.length-1]}
                                                    onChange={(event) => handleContainerChange(event, i)}
                                                >
                                                    {row.containers.map((container: string, index: number) => (
                                                        <MenuItem key={index} value={container}>
                                                            {container}
                                                        </MenuItem>
                                                    ))}
                                                </Select>
                                            </FormControl>
                                        )}
                                    </TableCell>
                                    <TableCell align="right">
                                        <Button onClick={() => onClickStreamLogs(row.podName, (selectedContainers[i] || row.containers[row.containers.length-1]))} variant="contained">Get Logs</Button>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}
            <Box mt={4}>
                <PodLogStreamClusterPage code={code} setCode={setCode} />
            </Box>
        </div>
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
    containers: string[];
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
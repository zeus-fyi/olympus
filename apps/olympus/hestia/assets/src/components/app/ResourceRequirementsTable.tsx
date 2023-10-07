import * as React from "react";
import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {TableContainer, TableFooter, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import {RootState} from "../../redux/store";
import {Cluster} from "../../redux/clusters/clusters.types";

export function ResourceRequirementsTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const resourceRequirements = createResourceRequirementsData(cluster);
    const dispatch = useDispatch();

    useEffect(() => {
        async function fetchData() {
            try {
            } catch (e) {
            }
        }
        fetchData();
    }, [dispatch]);

    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };
    const handleChangePage = (
        event: React.MouseEvent<HTMLButtonElement> | null,
        newPage: number,
    ) => {
        setPage(newPage);
    };

    // const handleClick = async (event: any, app: any) => {
    //     event.preventDefault();
    //     navigate('/app/'+app.topologySystemComponentID);
    // }

    if (cluster == null) {
        return (<div></div>)
    }

    if (resourceRequirements === undefined || resourceRequirements === null) {
        return (<div></div>)
    }
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - resourceRequirements.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 400 }} aria-label="app resource requirements table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >ClusterBase</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Workload</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">vCPU</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Memory</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Disk</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Count</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {resourceRequirements.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.componentBaseName}
                            </TableCell>
                            <TableCell align="left">{row.skeletonBaseName}</TableCell>
                            <TableCell align="left">{row.resourceSumsCPU === '0' ? '-' : row.resourceSumsCPU}</TableCell>
                            <TableCell align="left">{row.resourceSumsMemory === '0' ? '-' : row.resourceSumsMemory}</TableCell>
                            <TableCell align="left">{row.resourceSumsDisk === '0' ? '-' : row.resourceSumsDisk}</TableCell>
                            <TableCell align="left">{row.replicas === '0' ? '-' : row.replicas + 'x'}</TableCell>
                        </TableRow>
                    ))}
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={4} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}

function createResourceRequirementsData(cluster: Cluster): Array<{componentBaseName: string, skeletonBaseName: string, resourceSumsCPU: string, resourceSumsMemory: string, resourceSumsDisk: string}> {
    const resourceRequirementsData = [];

    for (const [componentBaseName, skeletonBases] of Object.entries(cluster.componentBases)) {
        for (const [skeletonBaseName, skeletonBase] of Object.entries(skeletonBases)) {
            if (skeletonBase.resourceSums) {
                const {cpuRequests, memRequests, diskRequests, replicas} = skeletonBase.resourceSums;

                if ((cpuRequests && cpuRequests !== '0') ||  (memRequests && memRequests !== '0') || (diskRequests && diskRequests !== '0')) {
                    resourceRequirementsData.push({
                        componentBaseName,
                        skeletonBaseName,
                        resourceSumsCPU: cpuRequests?.toString() ?? '',
                        resourceSumsMemory: memRequests?.toString() ?? '',
                        resourceSumsDisk: diskRequests?.toString() ?? '',
                        replicas: replicas?.toString() ?? '',
                    });
                }
            }
        }
    }

    return resourceRequirementsData;
}

export function createDiskResourceRequirements(cluster: Cluster): Array<{componentBaseName: string, skeletonBaseName: string, resourceSumsDisk: string, replicas: string, blockStorageCostUnit: number}> {
    const resourceRequirementsData = [];
    for (const [componentBaseName, skeletonBases] of Object.entries(cluster.componentBases)) {
        for (const [skeletonBaseName, skeletonBase] of Object.entries(skeletonBases)) {
            if (skeletonBase.resourceSums) {
                const {cpuRequests, memRequests, diskRequests, replicas} = skeletonBase.resourceSums;
                let blockStorageCostUnit = divideBy100GiB(diskRequests.toString());
                if ((cpuRequests && cpuRequests !== '0') ||  (memRequests && memRequests !== '0') || (diskRequests && diskRequests !== '0')) {
                    resourceRequirementsData.push({
                        componentBaseName,
                        skeletonBaseName,
                        resourceSumsDisk: diskRequests?.toString() ?? '',
                        replicas: replicas?.toString() ?? '',
                        blockStorageCostUnit: blockStorageCostUnit,
                    });

                }
            }
        }
    }
    return resourceRequirementsData;
}
type DiskSize = {
    value: number;
    unit?: 'B' | 'KiB' | 'Ki' | 'MiB' | 'Mi' | 'GiB' | 'Gi' | 'TiB' | 'Ti' | 'PiB' | 'Pi' | 'EiB' | 'Ei';
};

function parseDiskSize(input: string): DiskSize {
    const regex = /^(\d+(?:\.\d+)?)\s*(([KMGTPE]i)?B?)$/i;
    const match = input.match(regex);

    if (!match) {
        return {value: 0};
    }

    const value = parseFloat(match[1]);
    const unit = (match[3]?.toUpperCase() + (match[4]?.toUpperCase() ?? "")) as 'B' | 'KiB' | 'Ki' | 'MiB' | 'Mi' | 'GiB' | 'Gi' | 'TiB' | 'Ti' | 'PiB' | 'Pi' | 'EiB' | 'Ei' | undefined;

    return { value, unit: unit && (unit.charAt(0) + unit.slice(1).toLowerCase()) as 'B' | 'KiB' | 'Ki' | 'MiB' | 'Mi' | 'GiB' | 'Gi' | 'TiB' | 'Ti' | 'PiB' | 'Pi' | 'EiB' | 'Ei' };
}

function convertToBibiBytes(size: DiskSize): number {
    if (size.value === 0) {
        return 0;
    }
    if (!size.unit) {
        return size.value;
    }

    const unitMap = {
        B: 1,
        KiB: 1 << 10,
        Ki: 1 << 10,
        MiB: 1 << 20,
        Mi: 1 << 20,
        GiB: 1 << 30,
        Gi: 1 << 30,
        TiB: BigInt(1) << BigInt(40),
        Ti: BigInt(1) << BigInt(40),
        PiB: BigInt(1) << BigInt(50),
        Pi: BigInt(1) << BigInt(50),
        EiB: BigInt(1) << BigInt(60),
        Ei: BigInt(1) << BigInt(60),
    };
    const unitMultiplier = unitMap[size.unit];
    return Number(BigInt(size.value) * BigInt(unitMultiplier));
}

export function convertToPercentage(allocatable: string, capacity: string): number {
    const diskSize = parseDiskSize(allocatable);
    const bibiBytes = convertToBibiBytes(diskSize);

    const diskSizeCap = parseDiskSize(capacity);
    const bibiBytesCap = convertToBibiBytes(diskSizeCap);
    return 100-100*Number(Number(BigInt(bibiBytes)) / Number(BigInt(bibiBytesCap)));
}

export function convertToMi(input: string): number {
    const diskSize = parseDiskSize(input);
    const bibiBytes = convertToBibiBytes(diskSize);
    return Number(BigInt(bibiBytes) / BigInt(1024**2));
}

export function divideBy100GiB(input: string): number {
    const diskSize = parseDiskSize(input);
    const bibiBytes = convertToBibiBytes(diskSize);
    return Number(BigInt(bibiBytes) / BigInt(100*1024**3));
}

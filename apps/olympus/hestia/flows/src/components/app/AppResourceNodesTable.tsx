import * as React from "react";
import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {RootState} from "../../redux/store";
import {resourcesApiGateway} from "../../gateway/resources";
import {setAppNodes} from "../../redux/resources/resources.reducer";
import {NodeAudit} from "../../redux/resources/resources.types";
import {convertToMi} from "./ResourceRequirementsTable";

export function AppResourceNodesResourcesTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const resources = useSelector((state: RootState) => state.resources.resources);
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const appNodes = useSelector((state: RootState) => state.resources.appNodes);

    useEffect(() => {
        async function fetchData() {
            try {
                if (cluster === undefined) {
                    return;
                }
                if (cluster.clusterName === '') {
                    return;
                }
                const response = await resourcesApiGateway.getAppResources(cluster);
                // console.log(response.data)
                const nodes = await response.data as NodeAudit[];
                dispatch(setAppNodes(nodes));
            } catch (e) {
                const nodes = [] as NodeAudit[];
            }
        }
        fetchData().then(r => {
        });
    }, [cluster,dispatch]);

    if (appNodes === undefined) {
        return (<div></div>)
    }
    if (appNodes.length === 0) {
        return (<div></div>)
    }
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


    if (appNodes.length === 0) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - appNodes.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 1000 }} aria-label="private apps pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >CloudProvider</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Region</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Slug</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Kubernetes Version</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >CPU Allocatable/Capacity</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >RAM Allocatable/Capacity</TableCell>
                        {/*<TableCell style={{ fontWeight: 'normal', color: 'white'}} >RAM Utilization %</TableCell>*/}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {appNodes && appNodes.map((row: any, i: number) => (
                        <TableRow
                            key={i}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.cloudProvider}
                            </TableCell>
                            <TableCell align="left">{row.region}</TableCell>
                            <TableCell align="left">{row.slug}</TableCell>
                            <TableCell align="left">{row.kubernetesVersion}</TableCell>
                            <TableCell align="left">{row.status.allocatable['cpu'] +' / ' + row.status.capacity['cpu'] + ' vCPUs'}</TableCell>
                            <TableCell align="left">{convertToMi(row.status.allocatable['memory'])+'Mi' + ' / ' + convertToMi(row.status.capacity['memory'])+'Mi'}</TableCell>
                            {/*<TableCell align="left">{convertToPercentage(row.status.allocatable['memory'], row.status.capacity['memory']).toFixed(2)}</TableCell>*/}
                        </TableRow>
                    ))}
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={4} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={[10, 25, 100, { label: 'All', value: -1 }]}
                            colSpan={4}
                            count={appNodes.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            SelectProps={{
                                inputProps: {
                                    'aria-label': 'rows per page',
                                },
                                native: true,
                            }}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            ActionsComponent={TablePaginationActions}
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}

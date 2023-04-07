import * as React from "react";
import {useEffect} from "react";
import {useDispatch} from "react-redux";
import {useNavigate} from "react-router-dom";
import {appsApiGateway} from "../../gateway/apps";
import {setPrivateOrgApps} from "../../redux/apps/apps.reducer";
import {TableContainer, TableFooter, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

export function ResourceRequirementsTable(props: any) {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const resourceRequirements = [{}];
    const dispatch = useDispatch();
    let navigate = useNavigate();

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await appsApiGateway.getPrivateApps();
                dispatch(setPrivateOrgApps(response));
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
    //     navigate('/clusters/app/'+app.topologySystemComponentID);
    // }

    if (resourceRequirements == null) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - resourceRequirements.length) : 0;

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 400 }} aria-label="app resource requirements table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#8991B0'}} >
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >ClusterBase</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Workload</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">CPU</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Memory</TableCell>
                        <TableCell style={{ color: 'white'}} align="left">Disk</TableCell>
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
                            <TableCell align="left">{row.cpu}</TableCell>
                            <TableCell align="left">{row.memory}</TableCell>
                            <TableCell align="left">{row.disk}</TableCell>
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

import * as React from "react";
import {useEffect} from "react";
import {TableContainer, TableFooter, TablePagination, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import TablePaginationActions from "@mui/material/TablePagination/TablePaginationActions";
import {useDispatch, useSelector} from "react-redux";
import {UserStatsRow} from "./UserStatsRow";
import {UserFlowStats} from "../../../../redux/flows/flows.actions";
import {setOpenAdminUserRow, setUserFlowStats} from "../../../../redux/flows/flows.reducer";
import {aiApiGateway} from "../../../../gateway/ai";

export function UserStatsTable(props: any) {
    const { csvExport, isAdminPanel } = props;
    const [page, setPage] = React.useState(0);
    const openRunsRow = useSelector((state: any) => state.ai.openRunsRow);
    // const selectedRuns = useSelector((state: any) => state.ai.selectedRuns);
    const isInternal = useSelector((state: any) => state.sessionState.isInternal);
    const openAdminUserRow = useSelector((state: any) => state.ai.openAdminUserRow);
    const userFlowStats = useSelector((state: any) => state.flows.userFlowStats);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [loading, setIsLoading] = React.useState(false);
    const dispatch = useDispatch();

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                setIsLoading(true); // Set loading to true
                const response = await aiApiGateway.getAdminDashboardStats();
                dispatch(setUserFlowStats(response.data));
            } catch (error) {
                // console.log("error", error);
            } finally {
                setIsLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData({});
    }, []);


    const handleOpen = (index: number) => {
        const isOpen = openRunsRow[index] || false;
        // if (!isOpen) {
        //     getDetails(index);
        // }
        dispatch(setOpenAdminUserRow({ ...openRunsRow, [index]: !isOpen }));
    };
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

    if (loading) {
        return (<div></div>)
    }
    if (userFlowStats === null || userFlowStats === undefined) {
        return (<div></div>)
    }

    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - userFlowStats.length) : 0;

    // const handleClick = (name: string) => {
    //     const currentIndex = selectedRuns.indexOf(name);
    //     const newSelected = [...selectedRuns];
    //
    //     if (currentIndex === -1) {
    //         newSelected.push(name);
    //     } else {
    //         newSelected.splice(currentIndex, 1);
    //     }
    //     dispatch(setSelectedRuns(newSelected));
    // };
    //
    // const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
    //     if (event.target.checked) {
    //         const newSelected = userFlowStats.map((wf: any, ind: number) => ind);
    //         dispatch(setSelectedRuns(newSelected));
    //         return;
    //     }
    //     dispatch(setSelectedRuns([]));
    // };
    return (
        <TableContainer sx={{ minWidth: 1800 }} component={Paper}>
            <Table aria-label="workflow run analysis pagination table">
                <TableHead>
                    <TableRow style={{ backgroundColor: '#333'}} >
                        {/*<TableCell padding="checkbox">*/}
                        {/*    <Checkbox*/}
                        {/*        color="primary"*/}
                        {/*        indeterminate={userFlowStats.length > 0 && selectedRuns.length < userFlowStats.length && selectedRuns.length > 0}*/}
                        {/*        checked={userFlowStats.length > 0 && selectedRuns.length === userFlowStats.length}*/}
                        {/*        onChange={handleSelectAllClick}*/}
                        {/*    />*/}
                        {/*</TableCell>*/}
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >OrgID</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Creation Date</TableCell>
                        {/*<TableCell style={{ fontWeight: 'normal', color: 'white', minWidth: 50}}>Progress</TableCell>*/}
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Email</TableCell>
                        <TableCell style={{ fontWeight: 'normal', color: 'white'}} >Flow Count</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {rowsPerPage > 0 && userFlowStats && userFlowStats.map((row: UserFlowStats, index: number) => {
                        // Use orchestrationStrID to find details in orchDetails, or fallback to row if not found
                        // const detailedRow = orchDetails[row.orchestration.orchestrationStrID] || row;
                        return (
                            <UserStatsRow
                                key={index}
                                row={row} // Pass detailedRow, which may contain additional details
                                open={openRunsRow[index] || false}
                                handleOpen={() => handleOpen(index)}
                                index={index}
                                csvExport={csvExport}
                                // handleClick={handleClick}
                                // checked={selectedRuns.indexOf(index) >= 0}
                            />
                        );
                    })}
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
                            count={userFlowStats.length}
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

import * as React from 'react';
import {useEffect, useState} from 'react';
import {createTheme, styled, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, {AppBarProps as MuiAppBarProps} from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate, useParams} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import {Card, CardContent, FormControl, InputLabel, MenuItem, Select, Stack, Tab, Tabs} from "@mui/material";
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import MainListItems from "../dashboard/listItems";
import {RootState} from "../../redux/store";
import {IrisOrgGroupRoutesRequest, loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {setEndpoints, setGroupEndpoints, setTableMetrics} from "../../redux/loadbalancing/loadbalancing.reducer";
import TextField from "@mui/material/TextField";
import {PlanUsagePieCharts} from "./charts/pie/UsagePieChart";
import {MetricsChart, TableMetricsCharts} from "./charts/radar/MetricsCharts";
import {LoadBalancingRoutesTable} from "./tables/LoadBalancingRoutesTable";
import {LoadBalancingMetricsTable} from "./tables/MetricsTable";
import {LoadBalancingPriorityScoreMetricsTable} from "./tables/PriorityScoreMetricsTable";
import {ProceduresCatalogTable} from "./tables/ProceduresCatalogTable";

const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
}

export const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

export const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        '& .MuiDrawer-paper': {
            position: 'relative',
            whiteSpace: 'nowrap',
            width: drawerWidth,
            transition: theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
            boxSizing: 'border-box',
            ...(!open && {
                overflowX: 'hidden',
                transition: theme.transitions.create('width', {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.leavingScreen,
                }),
                width: theme.spacing(7),
                [theme.breakpoints.up('sm')]: {
                    width: theme.spacing(9),
                },
            }),
        },
    }),
);

const mdTheme = createTheme();

function LoadBalancingDashboardContent(props: any) {
    const params = useParams();
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }
    const endpoints = useSelector((state: RootState) => state.loadBalancing.routes);
    const groups = useSelector((state: RootState) => state.loadBalancing.groups);
    const [loading, setLoading] = useState(false);
    const [selected, setSelected] = useState<string[]>([]);
    const [groupName, setGroupName] = useState<string>("-all");
    const [tableRoutes, setTableRoutes] = useState<string[]>([]);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [page, setPage] = React.useState(0);
    const [isAdding, setIsAdding] = useState<boolean>(false);
    const [isAddingGroup, setIsAddingGroup] = useState<boolean>(false);
    const [isUpdatingGroup, setIsUpdatingGroup] = useState<boolean>(false);
    const [newEndpoint, setNewEndpoint] = useState<string>("");
    const [reload, setReload] = useState(false); // State to trigger reload
    const [createGroupName, setCreateGroupName] = React.useState("");
    const [selectedTab, setSelectedTab] = useState(0);
    const [selectedMainTab, setSelectedMainTab] = useState(0);

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                setLoading(true); // Set loading to true
                const response = await loadBalancingApiGateway.getEndpoints();
                dispatch(setEndpoints(response.data.routes));
                dispatch(setGroupEndpoints(response.data.orgGroupsRoutes));
                setTableRoutes(response.data.routes);
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData(params);
    }, [reload]);


    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true); // Set loading to true
                const response = await loadBalancingApiGateway.getTableMetrics(groupName);
                dispatch(setTableMetrics(response.data));
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData();
    }, [groupName]);

    const handleClick = (name: string) => {
        const currentIndex = selected.indexOf(name);
        const newSelected = [...selected];

        if (currentIndex === -1) {
            newSelected.push(name);
        } else {
            newSelected.splice(currentIndex, 1);
        }

        setSelected(newSelected);
    };

    const handleDeleteEndpointsSubmission = async () => {
        if (selected.length <= 0) {
            return;
        }
        try {
            setLoading(true); // Set loading to false regardless of success or failure.
            const selectedSet = new Set(selected); // Create a Set for O(1) lookup
            if (groupName === "-all"  || groupName === "unused") {
                const payload = {
                    routes: selected // Filter tableRoutes
                };
                const response = await loadBalancingApiGateway.deleteEndpoints(payload);
            } else {
                const payload: IrisOrgGroupRoutesRequest = {
                    groupName: groupName,
                    routes: tableRoutes.filter(route => !selectedSet.has(route)) // Filter tableRoutes
                };
                if (payload.routes.length === 0) {
                    const response = await loadBalancingApiGateway.deleteEndpoints(payload);
                } else {
                    const payloadPartial: IrisOrgGroupRoutesRequest = {
                        groupName: groupName,
                        routes: selected
                    }
                    const response = await loadBalancingApiGateway.removeEndpointsFromGroupRoutingTable(payloadPartial);
                }
            }
        } catch (error) {
            console.log("error", error);
        } finally {
            setLoading(false); // Set loading to false regardless of success or failure.
            setReload(!reload); // Trigger reload by flipping the state
        }
        setSelected([]);
    }

    const handleSubmitNewEndpointSubmission = async () => {
        if (newEndpoint) {
            setLoading(true); // Set loading to false regardless of success or failure.
            const payload: IrisOrgGroupRoutesRequest = {
                routes: [newEndpoint]
            };
            try {
                const response = await loadBalancingApiGateway.createEndpoints(payload);
                console.log(response.status)
                // handle the response accordingly
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
                setReload(!reload); // Trigger reload by flipping the state
            }
        }
        setIsAdding(false);
        setNewEndpoint("");
    };

    const handleSubmitNewGroupSubmission = async () => {
        if (selected.length > 0) {
            setLoading(true); // Set loading to false regardless of success or failure.
            const payload: IrisOrgGroupRoutesRequest = {
                groupName: createGroupName,
                routes: selected
            };
            try {
                const response = await loadBalancingApiGateway.createEndpoints(payload);
                console.log(response.status)
                // handle the response accordingly
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
                setReload(!reload); // Trigger reload by flipping the state
            }
        }
        setIsUpdatingGroup(false)
        setIsAddingGroup(false);
        setGroupName("-all");
    };

    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.checked) {
            const newSelected = tableRoutes.map((endpoint) => endpoint);
            setSelected(newSelected);
            return;
        }
        setSelected([]);
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

    const handleChangeGroup = (name: string) => {
        setPage(0);
        setSelected([]);
        setGroupName(name);
        setIsUpdatingGroup(false);
        if (endpoints == null || endpoints.length == 0) {
            setReload(!reload); // Trigger reload by flipping the state
        }
        if (name === "-all" || name === "unused") {
            setSelectedTab(0);
        }
        setTableRoutes(name === "-all" ? endpoints : groups[name]);
    };

    const handleUpdateGroupTableEndpointsSubmission = async () => {
        try {
            setLoading(true); // Set loading to false regardless of success or failure.
            const selectedSet = new Set(selected); // Create a Set for O(1) lookup
            const payload: IrisOrgGroupRoutesRequest = {
                groupName: groupName,
                routes: tableRoutes.filter(route => selectedSet.has(route)) // Filter tableRoutes
            };
            const response = await loadBalancingApiGateway.updateGroupRoutingTable(payload);
        } catch (error) {
            console.log("error", error);
        } finally {
            setLoading(false); // Set loading to false regardless of success or failure.
            setReload(!reload); // Trigger reload by flipping the state
        }
    }

    const handleAddGroupTableEndpointsSubmission = async () => {
        try {
            setLoading(true); // Set loading to true
            const selectedSet = new Set(selected); // Create a Set for O(1) lookup
            groups[groupName].forEach(route => selectedSet.add(route)); // Add each route from tableRoutes to selectedSet

            const payload: IrisOrgGroupRoutesRequest = {
                groupName: groupName,
                routes: Array.from(selectedSet) // Convert the Set back to an array
            };
            const response = await loadBalancingApiGateway.updateGroupRoutingTable(payload);
        } catch (error) {
            console.log("error", error);
        } finally {
            setLoading(false); // Set loading to false regardless of success or failure.
            setReload(!reload); // Trigger reload by flipping the state
        }
    }

    const handleClickAddGroupEndpoints = () => {
        setIsUpdatingGroup(true);
        setSelected([]);
        const filteredRoutes = endpoints.filter(
            endpoint => !groups[groupName].includes(endpoint)
        );

        setTableRoutes(filteredRoutes);
        return;
    };

    const handleClickViewGroupEndpoints = () => {
        setIsUpdatingGroup(false);
        setSelected([]);
        setTableRoutes(groups[groupName])
    };

    const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedTab(newValue);
    };

    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    };

    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open} style={{ backgroundColor: '#333'}}>
                    <Toolbar
                        sx={{
                            pr: '24px', // keep right padding when drawer closed
                        }}
                    >
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && { display: 'none' }),
                            }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{ flexGrow: 1 }}
                        >
                            Dashboard
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={handleLogout}
                        >Logout
                        </Button>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeftIcon />
                        </IconButton>
                    </Toolbar>
                    <Divider />
                    <List component="nav">
                        <MainListItems />
                        <Divider sx={{ my: 1 }} />
                    </List>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: (theme) =>
                            theme.palette.mode === 'light'
                                ? theme.palette.grey[100]
                                : theme.palette.grey[900],
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                    }}
                >
                    <Toolbar />
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction={"row"} spacing={2} >
                        <Card sx={{ maxWidth: 700 }}>
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="div">
                                    Load Balancing Management
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Summary of your routing groups & endpoints. Refresh the page if making updates to see the latest.
                                </Typography>
                            </CardContent>
                        <Box mr={2} ml={2} mt={2} mb={4}>
                            <Stack direction={"row"} spacing={2} alignItems={"center"}>
                                <FormControl sx={{  }} fullWidth variant="outlined">
                                    <InputLabel key={`groupNameLabel`} id={`groupName`}>
                                        Routing Group
                                    </InputLabel>
                                    <Select
                                        labelId={`groupNameLabel`}
                                        id={`groupName`}
                                        name="groupName"
                                        value={groupName}
                                        onChange={(event) => handleChangeGroup(event.target.value)}
                                        label="Routing Group"
                                    >
                                        <MenuItem key={'all'} value={'-all'}>{"all"}</MenuItem>
                                        {Object.keys(groups).map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                                    </Select>
                                </FormControl>
                            </Stack>
                        </Box>
                        <Box display="flex" m={2} mb={4}>
                            {groupName === "-all" &&
                                <Box mr={2}>
                                    <Button variant="contained" onClick={() => setIsAdding(!isAdding)}>
                                        Add Endpoints
                                    </Button>
                                </Box>
                            }
                            {groupName !== "-all" && groupName !== "unused" && !isUpdatingGroup &&
                                <Box mr={2}>
                                    <Button variant="contained" onClick={() => handleClickAddGroupEndpoints()}>
                                        Add Group Endpoints
                                    </Button>
                                </Box>
                            }
                            {groupName !== "-all" && groupName !== "unused" && isUpdatingGroup &&
                                <Box mr={2}>
                                    <Button variant="contained" onClick={() => handleClickViewGroupEndpoints()}>
                                        View Group Endpoints
                                    </Button>
                                </Box>
                            }
                            {groupName === "-all" &&
                                <Box>
                                    <Button variant="contained" onClick={() => setIsAddingGroup(!isAddingGroup)}>
                                        Add Groups
                                    </Button>
                                </Box>
                            }
                        </Box>
                            {isAddingGroup && groupName === "-all" && (
                                <Box m={2} mb={4} display="flex" alignItems="center">
                                    <Box flexGrow={1} mr={2}>
                                        <TextField
                                            fullWidth
                                            id="group-input"
                                            label="Group Name"
                                            variant="outlined"
                                            value={createGroupName}
                                            onChange={(e) => setCreateGroupName(e.target.value)}
                                        />
                                    </Box>
                                    <Box>
                                        <Button variant="contained" color="primary"  onClick={() => handleSubmitNewGroupSubmission()}>
                                            Submit
                                        </Button>
                                    </Box>
                                </Box>
                            )}
                        </Card>
                            {(groupName === "-all" || groupName === "unused") && (
                                <PlanUsagePieCharts reload={reload} setReload={setReload}/>
                            )}
                            {(groupName !== "-all" && groupName !== "unused") && (
                                <MetricsChart />
                            )}
                        </Stack>
                    </Container>
                    {groupName !== "-all" && groupName !== "unused" && (
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Box sx={{ mb: 2 }}>
                                <TableMetricsCharts tableName={groupName}/>
                            </Box>
                        </Container>
                    )}
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        {groupName !== "-all" && groupName !== "unused" && (
                            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                                <Tabs value={selectedTab} onChange={handleTabChange} aria-label="basic tabs example">
                                    <Tab label="Endpoints"  />
                                    <Tab label="Metrics"  />
                                    <Tab label="Priority Scores" />
                                    <Tab label="Procedures" />
                                </Tabs>
                            </Box>
                        )}
                        {groupName === "-all" && (
                            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                                <Tabs value={selectedTab} onChange={handleTabChange} aria-label="basic tabs example">
                                    <Tab label="Endpoints"  />
                                    <Tab label="Procedures" />
                                </Tabs>
                            </Box>
                        )}
                        {selectedTab === 0 && (
                        <LoadBalancingRoutesTable
                            selectedTab={selectedTab}
                            handleTabChange={handleTabChange}
                            page={page}
                            rowsPerPage={rowsPerPage}
                            loading={loading}
                            endpoints={tableRoutes}
                            groups={groups}
                            groupName={groupName}
                            selected={selected}
                            handleSelectAllClick={handleSelectAllClick}
                            handleClick={handleClick}
                            handleChangeRowsPerPage={handleChangeRowsPerPage}
                            handleChangePage={handleChangePage}
                            isAdding={isAdding}
                            setIsAdding={setIsAdding}
                            newEndpoint={newEndpoint}
                            isUpdatingGroup={isUpdatingGroup}
                            setNewEndpoint={setNewEndpoint}
                            handleSubmitNewEndpointSubmission={handleSubmitNewEndpointSubmission}
                            handleDeleteEndpointsSubmission={handleDeleteEndpointsSubmission}
                            handleUpdateGroupTableEndpointsSubmission={handleUpdateGroupTableEndpointsSubmission}
                            handleAddGroupTableEndpointsSubmission={handleAddGroupTableEndpointsSubmission}
                        />)}
                        {selectedTab === 1 && groupName !== "-all" && groupName !== "unused" &&  (
                            <LoadBalancingMetricsTable
                                selectedTab={selectedTab}
                                handleTabChange={handleTabChange}
                                page={page}
                                rowsPerPage={rowsPerPage}
                                loading={loading}
                                endpoints={tableRoutes}
                                groups={groups}
                                groupName={groupName}
                                selected={selected}
                                handleSelectAllClick={handleSelectAllClick}
                                handleClick={handleClick}
                                handleChangeRowsPerPage={handleChangeRowsPerPage}
                                handleChangePage={handleChangePage}
                                isAdding={isAdding}
                                setIsAdding={setIsAdding}
                                newEndpoint={newEndpoint}
                                isUpdatingGroup={isUpdatingGroup}
                                setNewEndpoint={setNewEndpoint}
                                handleSubmitNewEndpointSubmission={handleSubmitNewEndpointSubmission}
                                handleDeleteEndpointsSubmission={handleDeleteEndpointsSubmission}
                                handleUpdateGroupTableEndpointsSubmission={handleUpdateGroupTableEndpointsSubmission}
                                handleAddGroupTableEndpointsSubmission={handleAddGroupTableEndpointsSubmission}
                            />)}
                        {selectedTab === 2 && groupName !== "-all" && groupName !== "unused" && (
                            <LoadBalancingPriorityScoreMetricsTable
                                selectedTab={selectedTab}
                                handleTabChange={handleTabChange}
                                page={page}
                                rowsPerPage={rowsPerPage}
                                loading={loading}
                                endpoints={tableRoutes}
                                groups={groups}
                                groupName={groupName}
                                selected={selected}
                                handleSelectAllClick={handleSelectAllClick}
                                handleClick={handleClick}
                                handleChangeRowsPerPage={handleChangeRowsPerPage}
                                handleChangePage={handleChangePage}
                                isAdding={isAdding}
                                setIsAdding={setIsAdding}
                                newEndpoint={newEndpoint}
                                isUpdatingGroup={isUpdatingGroup}
                                setNewEndpoint={setNewEndpoint}
                                handleSubmitNewEndpointSubmission={handleSubmitNewEndpointSubmission}
                                handleDeleteEndpointsSubmission={handleDeleteEndpointsSubmission}
                                handleUpdateGroupTableEndpointsSubmission={handleUpdateGroupTableEndpointsSubmission}
                                handleAddGroupTableEndpointsSubmission={handleAddGroupTableEndpointsSubmission}
                            />)}
                        {selectedTab === 3 && groupName !== "-all" && groupName !== "unused" && (
                            <ProceduresCatalogTable
                                selectedTab={selectedTab}
                                handleTabChange={handleTabChange}
                                page={page}
                                rowsPerPage={rowsPerPage}
                                loading={loading}
                                endpoints={tableRoutes}
                                groups={groups}
                                groupName={groupName}
                                selected={selected}
                                handleSelectAllClick={handleSelectAllClick}
                                handleClick={handleClick}
                                handleChangeRowsPerPage={handleChangeRowsPerPage}
                                handleChangePage={handleChangePage}
                                isAdding={isAdding}
                                setIsAdding={setIsAdding}
                                newEndpoint={newEndpoint}
                                isUpdatingGroup={isUpdatingGroup}
                                setNewEndpoint={setNewEndpoint}
                                handleSubmitNewEndpointSubmission={handleSubmitNewEndpointSubmission}
                                handleDeleteEndpointsSubmission={handleDeleteEndpointsSubmission}
                                handleUpdateGroupTableEndpointsSubmission={handleUpdateGroupTableEndpointsSubmission}
                                handleAddGroupTableEndpointsSubmission={handleAddGroupTableEndpointsSubmission}
                            />)}
                    </Container>
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function LoadBalancingDashboard() {
    return <LoadBalancingDashboardContent />;
}
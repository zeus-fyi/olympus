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
import {
    Card,
    CardContent,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Select,
    Slider,
    Stack,
    Switch,
    Tab,
    Tabs
} from "@mui/material";
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import MainListItems from "../dashboard/listItems";
import {RootState} from "../../redux/store";
import {IrisOrgGroupRoutesRequest, loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {
    setEndpoints,
    setGroupEndpoints,
    setTableMetrics,
    setUserPlanDetails,
} from "../../redux/loadbalancing/loadbalancing.reducer";
import TextField from "@mui/material/TextField";
import {PlanUsagePieCharts} from "./charts/pie/UsagePieChart";
import {MetricsChart, TableMetricsCharts} from "./charts/radar/MetricsCharts";
import {LoadBalancingRoutesTable} from "./tables/LoadBalancingRoutesTable";
import {LoadBalancingMetricsTable} from "./tables/MetricsTable";
import {LoadBalancingPriorityScoreMetricsTable} from "./tables/PriorityScoreMetricsTable";
import {ProceduresCatalogTable} from "./tables/ProceduresCatalogTable";
import {IrisApiGateway} from "../../gateway/iris";
import JoyrideTutorialBegin, {State} from "./joyride/Joyride";
import {CallBackProps, STATUS} from "react-joyride";
import {useSetState} from "react-use";
import Checkbox from "@mui/material/Checkbox";
import {findKeyWithPrefix} from "./markdown/ExampleRequests";

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
    const planDetails = useSelector((state: RootState) => state.loadBalancing.planUsageDetails);
    const [planName, setPlanName] = useState<string>(planDetails?.planName.toLowerCase() ?? "standard");
    const [runTutorial, setRunTutorial] = useState<boolean>(planDetails?.tableUsage?.tutorialOn ?? true);
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
    const [tabCount, setTabCount] = useState(2);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const tableMetrics = useSelector((state: RootState) => state.loadBalancing.tableMetrics);
    const [loadingMetrics, setLoadingMetrics] = React.useState(false);
    const [sliderLatencyValue, setSliderLatencyValue] = useState( tableMetrics?.scaleFactors?.latencyScaleFactor ?? 0.6);
    const [sliderErrorValue, setSliderErrorValue] = useState(tableMetrics?.scaleFactors?.errorScaleFactor ?? 3.0);
    const [sliderDecayValue, setSliderDecayValue] = useState(tableMetrics?.scaleFactors?.decayScaleFactor ?? 0.95);

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
                setLoadingMetrics(true); // Set loading to true
                const response = await loadBalancingApiGateway.getTableMetrics(groupName);
                if (response.data === null) {
                    return;
                }
                dispatch(setTableMetrics(response.data));
                setSliderLatencyValue(response.data.scaleFactors.latencyScaleFactor);
                setSliderErrorValue(response.data.scaleFactors.errorScaleFactor);
                setSliderDecayValue(response.data.scaleFactors.decayScaleFactor);
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false)
                setLoadingMetrics(false); // Set loading to false regardless of success or failure.
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
            const selectedSet = new Set(selected); // Create a Set
            // for O(1) lookup
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

        if (name === "-all" || name === "unused") {
            setSelectedTab(0);
        } else {
            setSelectedMainTab(0)
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
    const handleScaleFactorChange = async (sf: string) => {
        let value: number;
        switch (sf) {
            case "latency":
                value = sliderLatencyValue;
                break;
            case "error":
                value = sliderErrorValue;
                break;
            case "decay":
                value = sliderDecayValue;
                break;
            default:
                return;
        }
        try {
            setLoadingMetrics(true); // Set loading to true
            const response = await IrisApiGateway.updateTableScaleFactor(groupName, sf, value);
        } catch (error) {
            console.log("error", error);
        } finally {
            setLoadingMetrics(false); // Set loading to false regardless of success or failure.
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

    // Handler for updating the slider value
    const onChangeLatencySlider = (event: any, newValue: number) => {
        setSliderLatencyValue(newValue);
    };
    // Handler for the "Set Default" button
    const handleSetDefaultLatency = () => {
        setSliderLatencyValue(0.6); // or some other default value
    };

    const onChangeErrorSlider = (event: any, newValue: number) => {
        setSliderErrorValue(newValue);
    };
    const handleSetDefaultError = () => {
        setSliderErrorValue(3.0); // or some other default value
    };

    const onChangeDecaySlider = (event: any, newValue: number) => {
        setSliderDecayValue(newValue);
    };

    const handleSetDefaultDecay= () => {
        setSliderDecayValue(0.95); // or some other default value
    };

    const onToggleTutorialSetting = async () => {
        try {
            const response = await loadBalancingApiGateway.updateTutorialSetting();
            if (response.data === null) {
                return;
            }
            let tmp = JSON.parse(JSON.stringify(planDetails));
            tmp.tableUsage.tutorialOn = response.data;
            dispatch(setUserPlanDetails(tmp));
            setRunTutorial(tmp.tableUsage.tutorialOn);
        } catch (error) {
            console.log("error", error);
        } finally {
        }
    }

    function CustomContent() {
        return (
            <div>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Load Balancing Dashboard: How to Get Started
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            <FormControlLabel
                                control={<Checkbox color="primary" onChange={onToggleTutorialSetting} />}
                                label="Disable and don't show this again"
                                style={{ fontSize: '1.2em' }}
                            />
                        </Typography>
                    </CardContent>
            </div>
        );
    }

    const allSteps = [
        {
            content: <CustomContent />,
            locale: { skip: <strong aria-label="skip"></strong> },
            placement: 'center',
            target: 'body',
        },
        {
            content: 'This view shows all the routes you have registered for use with the Load Balancer.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-all-routes',
            title: 'All Routes',
        },
        {
            content: 'This read-only view shows all registered routing procedures, only standard plans or higher can use procedures.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-all-procedures',
            title: 'All Procedures',
        },
        {
            content: 'Select a routing table to toggle the table view.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-qn-routing-table',
            title: 'QuickNode Generated Routing Table',
        },
        {
            content: 'Click on view details for your matching protocol and send request.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-procedures',
            title: 'Procedure Demo',
        },
        {
            content: 'You\'ll see the metrics from this request shortly after the request is sent.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-procedures',
            title: 'Procedure Demo',
        },
        {
            content: 'This view shows your priority score routes table & scale factors.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-priority-scores',
            title: 'Priority Scores',
        },
        {
            content: 'This view shows your route metrics.',
            placement: 'bottom',
            target: '.onboarding-card-highlight-metrics',
            title: 'Metrics',
        },
    ];

    const stepsForPlan = planName.toLowerCase() === 'lite' ? allSteps.slice(0, 4) : allSteps;
    const [{ run, steps }, setState] = useSetState<State>({
        run: runTutorial,
        // @ts-ignore
        steps: stepsForPlan,
    });

    const createJoyrideCallback = (plan: string) => (data: CallBackProps) => {
        const { status, index } = data;
        const finishedStatuses: string[] = [STATUS.FINISHED, STATUS.SKIPPED];

        if (status === STATUS.RUNNING) {
            if (plan.toLowerCase() === 'lite' && index > 3) {
                setState({ run: false });
                return;
            }
            switch (index) {
                case 0:
                    setSelectedMainTab(0);
                    setSelectedTab(0);
                    const gnAll = '-all'
                    setGroupName((gnAll));
                    setTableRoutes(endpoints);
                    break;
                case 2:
                    setSelectedMainTab(1);
                    break;
                case 3:
                    setSelectedMainTab(0);
                    setSelectedTab(0);
                    const gn = findKeyWithPrefix(groups);
                    setGroupName((gn));
                    setTableRoutes(groups[gn]);
                    break;
                case 4:
                    setSelectedMainTab(0);
                    setSelectedTab(3);
                    break;
                case 5:
                    setSelectedTab(3);
                    break;
                case 6:
                    setSelectedTab(2);
                    break;
                case 7:
                    setSelectedTab(1);
                    break;
                default:
                    break;
            }
        }

        if (finishedStatuses.includes(status)) {
            setState({ run: false });
        }
    };

    if (loading) {
        return <div></div>
    }
    if (loadingMetrics) {
        return <div></div>
    }

    return (
        <ThemeProvider theme={mdTheme}>
            <JoyrideTutorialBegin handleChangeGroup={handleChangeGroup}
                                  run={runTutorial}
                                  steps={steps}
                                  groups={groups}
                                  groupName={groupName}
                                  handleJoyrideCallback={createJoyrideCallback(planName.toLowerCase())}
            />
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
                        <Card  className="onboarding-card-highlight-qn-routing-table"  sx={{ maxWidth: 700 }}>
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
                                <Tabs value={selectedTab} onChange={handleTabChange} aria-label="basic tabs">
                                    <Tab label="Routes"  />
                                    {planName.toLowerCase() !== "lite" && (<Tab label="Metrics" className="onboarding-card-highlight-metrics" />)}
                                    {planName.toLowerCase() !== "lite" && (<Tab label="Priority Scores" className="onboarding-card-highlight-priority-scores"/>)}
                                    {planName.toLowerCase() !== "lite" && (<Tab className="onboarding-card-highlight-procedures" label="Procedures" />)}
                                </Tabs>
                            </Box>
                        )}
                        {groupName === "-all" && (
                            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                                <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                                    <Tab className="onboarding-card-highlight-all-routes" label="Routes"  />
                                    <Tab className="onboarding-card-highlight-all-procedures" label="Procedures" />
                                    <Tab label="Settings" />
                                </Tabs>
                            </Box>
                        )}
                        { (selectedTab === 0 && selectedMainTab === 0) && (
                        <LoadBalancingRoutesTable
                            selectedMainTab={selectedMainTab}
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
                                loadingMetrics={loadingMetrics}
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
                            <div>
                                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                    <Stack direction={"row"}>
                                    <Card sx={{ maxWidth: 700, minHeight: 125, mr: 2}}>
                                        <Stack direction={"column"}>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Priority Score Latency Scale Factor
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Adjust the slider to change the latency scale factor for the priority score.
                                                Your adjusted priority score is calculated as newScore = currentScore x (latency(percentile) + latencyScaleFactor).
                                            </Typography>
                                        </CardContent>
                                        <Box
                                            display="flex"
                                            alignItems="center"
                                            justifyContent="center"
                                            height="100%"
                                        >
                                            <Slider
                                                sx={{ mt: 4, mb: 10, ml:4, mr:4 }}
                                                aria-label="Always visible"
                                                min={0}
                                                max={1}
                                                step={0.001}
                                                value={sliderLatencyValue}
                                                onChange={(e, newValue) => onChangeLatencySlider(e, newValue as number)}
                                                valueLabelDisplay="on"
                                            />
                                        </Box>
                                            <CardContent>
                                                <Stack direction={"row"} spacing={2}>
                                                <Button variant="contained" fullWidth color="primary" onClick={handleSetDefaultLatency}>
                                                    Set Default
                                                </Button>
                                                <Button variant="contained" fullWidth color="primary" onClick={() => handleScaleFactorChange("latency")}>
                                                    Update
                                                </Button>
                                                </Stack>
                                            </CardContent>
                                    </Stack>
                                    </Card>
                                    <Card sx={{ maxWidth: 700, minHeight: 125,  mr: 2}}>
                                        <Stack direction={"column"}>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Error Scale Factor
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Whenever a 4xx or 5xx error is returned, the error scale factor is used to adjust the priority score.
                                                Your adjusted priority score is calculated as newScore = currentScore x errorScaleFactor.
                                            </Typography>
                                        </CardContent>
                                        <Box
                                            display="flex"
                                            alignItems="center"
                                            justifyContent="center"
                                            height="100%"
                                        >
                                            <Slider
                                                sx={{ mt: 4, mb: 10, ml:4, mr:4 }}
                                                aria-label="Always visible"
                                                min={0}
                                                max={10}
                                                onChange={(e, newValue) => onChangeErrorSlider(e, newValue as number)}
                                                step={0.001}
                                                value={sliderErrorValue}
                                                valueLabelDisplay="on"
                                            />
                                        </Box>
                                        <CardContent>
                                            <Stack direction={"row"} spacing={2}>
                                                <Button variant="contained" fullWidth color="primary" onClick={handleSetDefaultError}>
                                                    Set Default
                                                </Button>
                                                <Button variant="contained" fullWidth color="primary" onClick={() => handleScaleFactorChange("error")}>
                                                    Update
                                                </Button>
                                            </Stack>
                                        </CardContent>
                                        </Stack>
                                    </Card>
                                    <Card sx={{ maxWidth: 700, minHeight: 125}}>
                                        <Stack direction={"column"}>
                                            <CardContent>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Decay Scale Factor
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    Whenever N (number of table endpoints == adaptive requests) have been made relative to an endpoint,
                                                    your adjusted priority score is calculated as newScore = currentScore x decayScaleFactor
                                                </Typography>
                                            </CardContent>
                                            <Box
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="center"
                                                height="100%"
                                            >
                                                <Slider
                                                    sx={{ mt: 4, mb: 10, ml:4, mr:4 }}
                                                    aria-label="Always visible"
                                                    min={0}
                                                    max={1}
                                                    step={0.001}
                                                    value={sliderDecayValue}
                                                    onChange={(e, newValue) => onChangeDecaySlider(e, newValue as number)}
                                                    valueLabelDisplay="on"
                                                />
                                            </Box>
                                            <CardContent>
                                                <Stack direction={"row"} spacing={2}>
                                                    <Button variant="contained" fullWidth color="primary" onClick={handleSetDefaultDecay}>
                                                        Set Default
                                                    </Button>
                                                    <Button variant="contained" fullWidth color="primary" onClick={() => handleScaleFactorChange("decay")}>
                                                        Update
                                                    </Button>
                                                </Stack>
                                            </CardContent>
                                        </Stack>
                                    </Card>
                                </Stack>
                                </Container>
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
                                />
                            </div>)}
                        {( selectedTab === tabCount +1 || selectedMainTab === 1 && groupName == "-all") && (
                                <ProceduresCatalogTable
                                    selectedMainTab={selectedMainTab}
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
                                />
                            )}
                        {( selectedMainTab === 2 && groupName == "-all") && (
                            <div>
                                <Box width="50%" sx={{ mt: 4, display: 'flex' }}>
                                    <Card >
                                        <CardContent>
                                            <Stack direction={"row"} spacing={2} alignItems="center">
                                                <Typography variant="body2" color="text.secondary" sx={{mt: 4}}>
                                                    Toggle on to re-enable the tutorial.
                                                </Typography>
                                                <Switch
                                                    sx={{ml: 2}}
                                                    checked={runTutorial}
                                                    onChange={onToggleTutorialSetting}
                                                    inputProps={{ 'aria-label': 'controlled' }}
                                                />
                                            </Stack>
                                        </CardContent>
                                    </Card>
                                </Box>
                            </div>
                        )}
                            </Container>
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function LoadBalancingDashboard(props: any) {
    return <LoadBalancingDashboardContent/>;
}
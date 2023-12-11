import * as React from 'react';
import {useCallback, useEffect, useState} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {
    Card,
    CardContent,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    Switch,
    Tab,
    Tabs,
    TextareaAutosize
} from "@mui/material";
import authProvider from "../../redux/auth/auth.actions";
import MainListItems from "../dashboard/listItems";
import {WorkflowTable} from "./WorkflowTable";
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import {AiSearchAnalysis} from "./AiAnalysisSummaries";
import TextField from "@mui/material/TextField";
import {AppBar, Drawer} from "../dashboard/Dashboard";
import {RootState} from "../../redux/store";
import {
    setDiscordOptionsCategoryName,
    setRetrievalKeywords,
    setRetrievalPlatform,
    setRetrievalPlatformGroups,
    setRetrievalPrompt,
    setRetrievalUsernames,
    setSearchIndexer,
    setSearchResults,
    setSelectedRuns,
    setSelectedSearchIndexers,
    setSelectedWorkflows,
    setWebRoutingGroup
} from "../../redux/ai/ai.reducer";
import {aiApiGateway} from "../../gateway/ai";
import {set} from 'date-fns';
import {TimeRange} from '@matiaslgonzalez/react-timeline-range-slider';
import {WorkflowAnalysisTable} from "./WorkflowAnalysisTable";
import {AiSearchParams, PostRunsActionRequest, PostWorkflowsActionRequest} from "../../redux/ai/ai.types";
import {setEndpoints, setGroupEndpoints} from "../../redux/loadbalancing/loadbalancing.reducer";
import {loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {SearchIndexersTable} from "./SearchIndexersTable";
import {isValidLabel} from "../clusters/wizard/builder/AddComponentBases";

const mdTheme = createTheme();
const analysisStart = "====================================================================================ANALYSIS====================================================================================\n"
const analysisDone = "====================================================================================ANALYSIS-DONE===============================================================================\n"

function AiWorkflowsDashboardContent(props: any) {
    const [open, setOpen] = useState(true);
    const [loading, setIsLoading] = useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const selected = useSelector((state: any) => state.ai.selectedWorkflows);
    const selectedRuns = useSelector((state: any) => state.ai.selectedRuns);
    const runs = useSelector((state: any) => state.ai.runs);
    const groups = useSelector((state: RootState) => state.loadBalancing.groups);
    const [code, setCode] = useState('');
    const [unixStartTime, setUnixStartTime] = useState(0);
    const [stepSize, setStepSize] = useState(1);
    const [stepSizeUnit, setStepSizeUnit] = useState('hours');
    const retrieval = useSelector((state: RootState) => state.ai.retrieval);
    const [analyzeNext, setAnalyzeNext] = useState(true);
    const [customBasePeriod, setCustomBasePeriod] = useState(false);
    const [customBasePeriodStepSize, setCustomBasePeriodStepSize] = useState(5);
    const [customBasePeriodStepSizeUnit, setCustomBasePeriodStepSizeUnit] = useState('minutes');
    const workflows = useSelector((state: any) => state.ai.workflows);
    const [requestRunsStatus, setRequestRunsStatus] = useState('');
    const [requestRunsStatusError, setRequestRunsStatusError] = useState('');
    const [requestStatus, setRequestStatus] = useState('');
    const [requestStatusError, setRequestStatusError] = useState('');
    const [requestIndexerStatus, setRequestIndexerStatus] = useState('');
    const [requestIndexerStatusError, setRequestIndexerStatusError] = useState('');
    const searchIndexer = useSelector((state: any) => state.ai.searchIndexer);
    const platformSecretReference = useSelector((state: any) => state.ai.platformSecretReference);
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                setIsLoading(true); // Set loading to true
                const response = await loadBalancingApiGateway.getEndpoints();
                dispatch(setEndpoints(response.data.routes));
                dispatch(setGroupEndpoints(response.data.orgGroupsRoutes));
            } catch (error) {
                console.log("error", error);
            } finally {
                setIsLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData({});
    }, []);
    const dispatch = useDispatch();
    const getCurrentUnixTimestamp = (): number => {
        return Math.floor(Date.now() / 1000);
    };
    const now = new Date();
    const getTodayAtSpecificHour = (hour: number = 12) =>
        set(now, { hours: hour, minutes: 0, seconds: 0, milliseconds: 0 });

    const [searchInterval, setSearchInterval] = useState<[Date, Date]>([getTodayAtSpecificHour(0), getTodayAtSpecificHour(24)]);
    const onTimeRangeChange = useCallback((interval: [Date, Date]) => {
        setSearchInterval(interval);
    }, []);

    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const handleToggleChange = (event: any) => {
        setAnalyzeNext(event.target.checked);
    };
    const handleToggleChangePeriod = (event: any) => {
        setCustomBasePeriod(event.target.checked);
    };

    const handleRunsActionRequest = async (event: any, action: string) => {
        const params: PostRunsActionRequest = {
            action: action,
            runs: selectedRuns.map((index: number) => {
                return runs[index].orchestration
            })
        }
        if (params.runs.length === 0) {
            return
        }
        try {
            setIsLoading(true)
            const response = await aiApiGateway.execRunsActionRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSelectedRuns([]));
                setRequestRunsStatus('Run cancellation submitted successfully')
                setRequestRunsStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestRunsStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestRunsStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    const handleWorkflowAction = async (event: any, action: string) => {
        const params: PostWorkflowsActionRequest = {
            action: action,
            unixStartTime: unixStartTime,
            duration: stepSize,
            durationUnit: stepSizeUnit,
            customBasePeriod: customBasePeriod,
            customBasePeriodStepSize: customBasePeriodStepSize,
            customBasePeriodStepSizeUnit: customBasePeriodStepSizeUnit,
            workflows: selected.map((index: number) => {
                return workflows[index]
            })
        }
        if (params.workflows.length === 0) {
            return
        }
        try {
            setIsLoading(true)
            const response = await aiApiGateway.execWorkflowsActionRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSelectedWorkflows([]));
                setRequestStatus('Workflow run started successfully')
                setRequestStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }
    function convertUnixToDateTimeLocal(unixTime: number): string {
        if (!unixTime) return '';
        const date = new Date(unixTime * 1000);
        const offset = date.getTimezoneOffset() * 60000; // offset in milliseconds
        return new Date(date.getTime() - offset).toISOString().slice(0, 16);
    }

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    const handleSearchRequest = async (timeRange: '1 hour'| '24 hours' | '7 days'| '30 days' | 'window' | 'all') => {
        try {
            setIsLoading(true)
            // Construct the retrieval object based on your existing variables
            const params: AiSearchParams = {
                timeRange,
                window: {
                    start: searchInterval[0],
                    end: searchInterval[1],
                    unixStartTime: searchInterval[0].getTime() / 1000, // converting milliseconds to seconds
                    unixEndTime: searchInterval[1].getTime() / 1000   // converting milliseconds to seconds
                },
                retrieval
            };
            const response = await aiApiGateway.searchRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSearchResults(data));
                setCode(data)
                setRequestStatus('')
                setRequestStatusError('')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    if (loading) {
        return <div>Loading...</div>;
    }

    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        if (newValue !== 1) {
            dispatch(setSelectedWorkflows([]))
        }
        if (newValue !== 3) {
            dispatch(setSelectedRuns([]))
        }
        setRequestStatus('')
        setRequestStatusError('')
        setRequestIndexerStatusError('')
        setRequestIndexerStatus('')
        setRequestRunsStatus('')
        setRequestRunsStatusError('')
        setSelectedMainTab(newValue);
    };

    const onChangeText = (textInput: string) => {
        setCode(textInput);
    };

    const formatTick2 = (ms: number) => {
        return new Date(ms).toLocaleTimeString([], { hour: 'numeric', hour12: true });
    };

    const handleSubmitIndexer = async (event: any) => {
        if (searchIndexer.platform.length === 0) {
            setRequestIndexerStatus('Search platform name cannot be empty.')
            setRequestIndexerStatusError('error')
            return
        }

        if (!isValidLabel(searchIndexer.searchGroupName)) {
            setRequestIndexerStatus('Search group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
            setRequestIndexerStatusError('error')
            return;
        }
        if (searchIndexer.query.length === 0) {
            setRequestIndexerStatus('Search indexer query required.');
            setRequestIndexerStatusError('error')
            return
        }
        try {
            setIsLoading(true)
            const params =  {
                searchIndexer,
                platformSecretReference
            }
            const response = await aiApiGateway.searchIndexerCreateOrUpdateActionRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSelectedSearchIndexers([]));
                setRequestIndexerStatus('Search indexer request sent successfully')
                setRequestIndexerStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestIndexerStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestIndexerStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    const AppBarAi = (props: any) => {
        const {toggleDrawer, open, handleLogout} = props;
        return (
            <div>
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
                        Mockingbird — A Time Series RAG-LLM Workflow Engine
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
            </div>
        )};

    const handleDateTimeChange = (e: any) => {
        const date = new Date(e.target.value);
        setUnixStartTime(date.getTime() / 1000); // Convert back to Unix timestamp
    };


    const ti = analyzeNext ? 'Next' : 'Previous';
    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBarAi toggleDrawer={toggleDrawer} open={open} handleLogout={handleLogout} />
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
                    { (selectedMainTab === 0) &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="row" spacing={2}>
                            <Card sx={{ minWidth: 100, maxWidth: 600 }}>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Search Previews
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        This allows you to search across platforms for data. Use this
                                        to help mock your AI workflow time series designs. You will need to either send us your data,
                                        or use our data retrieval indexer with your relevant platform API key & indexing query to get the data you want to analyze.
                                        Send us a message at support@zeus.fyi for indexer access.
                                    </Typography>
                                </CardContent>
                                <CardContent>
                                    <div>
                                        <Stack direction="column" spacing={2} sx={{ mt: 0, mb: 0 }}>
                                            {/*<Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>*/}
                                            {/*    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>*/}
                                            {/*        <TextField*/}
                                            {/*            fullWidth*/}
                                            {/*            id="retrieval-name"*/}
                                            {/*            label="Retrieval Name"*/}
                                            {/*            variant="outlined"*/}
                                            {/*            value={retrieval.retrievalName}*/}
                                            {/*            onChange={(e) => dispatch(setRetrievalName(e.target.value))}*/}
                                            {/*        />*/}
                                            {/*    </Box>*/}
                                            {/*    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>*/}
                                            {/*        <TextField*/}
                                            {/*            fullWidth*/}
                                            {/*            id="retrieval-group"*/}
                                            {/*            label="Retrieval Group"*/}
                                            {/*            variant="outlined"*/}
                                            {/*            value={retrieval.retrievalGroup}*/}
                                            {/*            onChange={(e) => dispatch(setRetrievalGroup(e.target.value))}*/}
                                            {/*        />*/}
                                            {/*    </Box>*/}
                                            {/*</Stack>*/}
                                            <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                <FormControl fullWidth>
                                                    <InputLabel id="platform-label">Platform</InputLabel>
                                                    <Select
                                                        labelId="platform-label"
                                                        id="platforms-input"
                                                        value={retrieval.retrievalPlatform}
                                                        label="Platform"
                                                        onChange={(e) => dispatch(setRetrievalPlatform(e.target.value))}
                                                    >
                                                        <MenuItem value="web">Web</MenuItem>
                                                        <MenuItem value="reddit">Reddit</MenuItem>
                                                        <MenuItem value="twitter">Twitter</MenuItem>
                                                        <MenuItem value="discord">Discord</MenuItem>
                                                        <MenuItem value="telegram">Telegram</MenuItem>
                                                    </Select>
                                                </FormControl>
                                            </Box>

                                            { retrieval.retrievalPlatform !== 'web' &&
                                                <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                    <TextField
                                                        fullWidth
                                                        id="group-input"
                                                        label={"Platform Groups"}
                                                        variant="outlined"
                                                        value={retrieval.retrievalPlatformGroups}
                                                        onChange={(e) => dispatch(setRetrievalPlatformGroups(e.target.value))}
                                                    />
                                                </Box>
                                            }
                                            { retrieval.retrievalPlatform === 'web' &&
                                                <div>
                                                        <Typography variant="h6" color="text.secondary">
                                                            Use a Load Balancer group for web data retrieval.
                                                        </Typography>
                                                        <FormControl sx={{ mt: 3 }} fullWidth variant="outlined">
                                                        <InputLabel key={`groupNameLabel`} id={`groupName`}>
                                                            Routing Group
                                                        </InputLabel>
                                                        <Select
                                                            labelId={`groupNameLabel`}
                                                            id={`groupName`}
                                                            name="groupName"
                                                            value={retrieval.webFilters?.routingGroup || ''}
                                                            onChange={(e) => dispatch(setWebRoutingGroup(e.target.value))}
                                                            label="Routing Group"
                                                        >
                                                            {Object.keys(groups).map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                                                        </Select>
                                                        </FormControl>
                                                </div>
                                            }
                                            { retrieval.retrievalPlatform === 'discord' &&
                                                <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                    <TextField
                                                        fullWidth
                                                        id="category-name-input"
                                                        label="Discord Category Name"
                                                        variant="outlined"
                                                        value={retrieval.discordFilters?.categoryName || ''}
                                                        onChange={(e) => dispatch(setDiscordOptionsCategoryName(e.target.value))}
                                                    />
                                                </Box>
                                            }
                                            { retrieval.retrievalPlatform !== 'web' &&
                                                <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                <TextField
                                                    fullWidth
                                                    id="usernames-input"
                                                    label="Usernames"
                                                    variant="outlined"
                                                    value={retrieval.retrievalUsernames}
                                                    onChange={(e) => dispatch(setRetrievalUsernames(e.target.value))}
                                                />
                                            </Box>
                                            }
                                            <Typography variant="h6" color="text.secondary">
                                                Search keywords using comma separated values below.
                                            </Typography>
                                            <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                <TextField
                                                    fullWidth
                                                    id="keywords-input"
                                                    label="Keywords"
                                                    variant="outlined"
                                                    value={retrieval.retrievalKeywords}
                                                    onChange={(e) => dispatch(setRetrievalKeywords(e.target.value))}
                                                />
                                            </Box>
                                            <Typography variant="h6" color="text.secondary">
                                                Optionally describe what you're looking for, and the AI will analyze your returned search data.
                                            </Typography>
                                            <Box  sx={{ mb: 2, mt: 2 }}>
                                                <TextareaAutosize
                                                    minRows={18}
                                                    value={retrieval.retrievalPrompt}
                                                    onChange={(e) => dispatch(setRetrievalPrompt(e.target.value))}
                                                    style={{ resize: "both", width: "100%" }}
                                                />
                                            </Box>
                                        </Stack>
                                    </div>
                                </CardContent>
                            </Card>
                            <Card sx={{ minWidth: 500, maxWidth: 900 }}>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Search Window
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Select a time window to search for data.
                                    </Typography>
                                    <Box flexGrow={1} sx={{ mt: 2, mb: -12 }}>
                                        <TimeRange
                                            ticksNumber={12}
                                            timelineInterval={[getTodayAtSpecificHour(0), getTodayAtSpecificHour(24)]}
                                            // @ts-ignore
                                            onChangeCallback={onTimeRangeChange}
                                            formatTick={formatTick2}
                                            showNow={true}
                                            selectedInterval={searchInterval}
                                        />
                                    </Box>
                                    <Button fullWidth variant="contained" onClick={() => handleSearchRequest('window')} >Search Window</Button>
                                    <Box flexGrow={1} sx={{ mt: 2 }}>
                                        <Typography variant="body2" color="text.secondary">
                                            Use these buttons to search previous time intervals relative to the current time.
                                        </Typography>
                                    </Box>
                                    <Box flexGrow={1} sx={{ mb: 2, mt: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchRequest('1 hour')} >Search 1 Hour</Button>
                                    </Box>
                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchRequest('24 hours')} >Search 24 Hours</Button>
                                    </Box>
                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchRequest('7 days')} >Search 7 Days</Button>
                                    </Box>
                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchRequest('30 days')} >Search 30 Days </Button>
                                    </Box>
                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchRequest('all')} >Search All Records</Button>
                                    </Box>
                                </CardContent>
                            </Card>
                        </Stack>
                    </Container>
                    }
                    { (selectedMainTab === 1) &&
                        <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
                            <Stack direction="row" spacing={2}>
                                <Card sx={{ minWidth: 100, maxWidth: 600 }}>
                                    <CardContent>
                                        <Typography gutterBottom variant="h5" component="div">
                                            Search Indexer
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary">
                                            This allows you to use our data retrieval indexer with your relevant platform API key & indexing query to get the data you want to analyze.
                                        </Typography>
                                    </CardContent>
                                    <CardContent>
                                        <div>
                                            <Stack direction="column" spacing={2} sx={{ mt: 0, mb: 0 }}>
                                                <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="platform-label">Platform</InputLabel>
                                                        <Select
                                                            labelId="indexer-platform-label"
                                                            id="indexer-platforms-input"
                                                            value={searchIndexer.platform}
                                                            label="Platform"
                                                            onChange={(e) => dispatch(setSearchIndexer({ ...searchIndexer, platform: e.target.value }))}
                                                        >
                                                            <MenuItem value="openai">OpenAI</MenuItem>
                                                            <MenuItem value="reddit">Reddit</MenuItem>
                                                            <MenuItem value="twitter">Twitter</MenuItem>
                                                            <MenuItem value="discord">Discord</MenuItem>
                                                            {/*<MenuItem value="telegram">Telegram</MenuItem>*/}
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                                { searchIndexer.platform !== 'openai' &&

                                                <div>
                                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                                        <TextField
                                                            fullWidth
                                                            id="search-group-input"
                                                            label={"Search Group Name"}
                                                            variant="outlined"
                                                            value={searchIndexer.searchGroupName}
                                                            onChange={(e) => dispatch(setSearchIndexer({ ...searchIndexer, searchGroupName: e.target.value }))}
                                                        />
                                                    </Box>
                                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                                        <TextField
                                                            fullWidth
                                                            id="search-group-input"
                                                            label={searchIndexer.platform === 'reddit' ? "Subreddits" : "Platform Search Query"}
                                                            variant="outlined"
                                                            value={searchIndexer.query}
                                                            onChange={(e) => dispatch(setSearchIndexer({ ...searchIndexer, query: e.target.value }))}
                                                        />
                                                    </Box>
                                                </div>
                                                }
                                                    <div>
                                                        <Typography variant="h6" color="text.secondary">
                                                           AWS Platform Secret Reference & Key For Search Indexer
                                                        </Typography>
                                                        { (searchIndexer.platform === 'discord' || searchIndexer.platform === 'openai') &&
                                                            <div>
                                                                <Stack direction={"row"} >
                                                                    <Box flexGrow={1} sx={{ mt: 2}}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-group-input"
                                                                            label={"Secret Name"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretKeyName === '' ? searchIndexer.platform + '-api-key' : platformSecretReference.secretKeyName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mt: 2, ml: 2 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-key-group-input"
                                                                            label={"Secret Key"}
                                                                            variant="outlined"

                                                                            value={platformSecretReference.secretGroupName}
                                                                        />
                                                                    </Box>
                                                                </Stack>
                                                            </div>
                                                        }
                                                        { (searchIndexer.platform === 'reddit' || searchIndexer.platform === 'twitter') &&
                                                            <div>
                                                                <Stack direction={"row"} >
                                                                    <Box flexGrow={1} sx={{ mt: 2}}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-group-input"
                                                                            label={"Secret Name"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretKeyName === '' ? searchIndexer.platform + '-oauth2-public' : platformSecretReference.secretKeyName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mt: 2, ml: 2 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-key-group-input"
                                                                            label={"Secret Key"}
                                                                            variant="outlined"

                                                                            value={platformSecretReference.secretGroupName}
                                                                        />
                                                                    </Box>
                                                                </Stack>
                                                                <Stack direction={"row"} >
                                                                    <Box flexGrow={1} sx={{ mt: 2}}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-group-input"
                                                                            label={"Secret Name"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretKeyName === '' ? searchIndexer.platform + '-oauth2-secret' : platformSecretReference.secretKeyName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mt: 2, ml: 2 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-key-group-input"
                                                                            label={"Secret Key"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretGroupName}
                                                                        />
                                                                    </Box>
                                                                </Stack>
                                                            </div>
                                                        }
                                                        { searchIndexer.platform === 'reddit' &&
                                                            <div >
                                                                <Stack direction={"row"} >
                                                                    <Box flexGrow={1} sx={{ mt: 2}}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-group-input"
                                                                            label={"Secret Name"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretKeyName === '' ? searchIndexer.platform + '-username' : platformSecretReference.secretKeyName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mt: 2, ml: 2 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-key-group-input"
                                                                            label={"Secret Key"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretGroupName}
                                                                        />
                                                                    </Box>
                                                                </Stack>
                                                                <Stack direction={"row"} >
                                                                    <Box flexGrow={1} sx={{ mt: 2}}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-group-input"
                                                                            label={"Secret Name"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretKeyName === '' ? searchIndexer.platform + '-password': platformSecretReference.secretKeyName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mt: 2, ml: 2 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id="platformSecretReference-key-group-input"
                                                                            label={"Secret Key"}
                                                                            variant="outlined"
                                                                            value={platformSecretReference.secretGroupName}
                                                                        />
                                                                    </Box>
                                                                </Stack>
                                                            </div>
                                                        }
                                                    </div>

                                                {requestIndexerStatus != '' && (
                                                    <Container sx={{  mt: 2}}>
                                                        <Typography variant="h6" color={requestIndexerStatusError}>
                                                            {requestIndexerStatus}
                                                        </Typography>
                                                    </Container>
                                                )}
                                                { searchIndexer.platform != 'openai' &&
                                                <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                                    <Button fullWidth variant="outlined"  onClick={(e) => handleSubmitIndexer(e)}>Start Indexing</Button>
                                                </Box>
                                                }
                                            </Stack>
                                        </div>
                                    </CardContent>
                                </Card>
                            </Stack>
                        </Container>
                    }
                    { (selectedMainTab === 2) &&
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Card sx={{ minWidth: 500, maxWidth: 1000 }}>
                            <CardContent>
                                <Box flexGrow={1} sx={{ mb: 2, mt: 0 }}>
                                    <Typography gutterBottom variant="h4" component="div">
                                        Run Scheduler
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Press Start to schedule a workflow that will run the analysis on the time intervals you've defined. It will
                                        process the data that gets generated from your search query, and then aggregate the results into over a rolling window.
                                    </Typography>
                                </Box>
                                    <Stack direction="row" spacing={2} sx={{ ml: 2, mr: 2, mt: 4, mb: 2 }}>
                                        <Box sx={{ width: '33%' }}> {/* Adjusted Box for TextField */}
                                            <TextField
                                                type="number"
                                                label="Runtime Duration"
                                                variant="outlined"
                                                inputProps={{ min: 1 }}  // Set minimum value to 1
                                                value={stepSize}
                                                onChange={(e)=> setStepSize(Number(e.target.value))}
                                                fullWidth
                                            />
                                        </Box>
                                        <Box sx={{ width: '33%' }}> {/* Adjusted Box for FormControl */}
                                            <FormControl fullWidth>
                                                <InputLabel id="time-unit-label">Time Unit</InputLabel>
                                                <Select
                                                    labelId="time-unit-label"
                                                    id="time-unit-select"
                                                    value={stepSizeUnit}
                                                    label="Time Unit"
                                                    onChange={(e) => setStepSizeUnit(e.target.value)}
                                                >
                                                    <MenuItem value="minutes">Minutes</MenuItem>
                                                    <MenuItem value="hours">Hours</MenuItem>
                                                    <MenuItem value="days">Days</MenuItem>
                                                    <MenuItem value="weeks">Weeks</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <FormControlLabel
                                            control={<Switch checked={analyzeNext} onChange={handleToggleChange} />}
                                            label={analyzeNext ? 'Analyze Next' : 'Analyze Previous'}
                                        />
                                    </Stack>
                                    { customBasePeriod &&
                                        <div>
                                            <Stack direction="row" spacing={2} sx={{ ml: 2, mr: 2, mt: 4, mb: 2 }}>
                                            <Box sx={{ width: '33%' }}>
                                                <TextField
                                                    type="number"
                                                    label="Override Workflow Step Size"
                                                    variant="outlined"
                                                    inputProps={{ min: 1 }}  // Set minimum value to 1
                                                    value={customBasePeriod ? customBasePeriodStepSize: 0}
                                                    onChange={(e) => setCustomBasePeriodStepSize(Number(e.target.value))}
                                                    fullWidth
                                                />
                                            </Box>
                                            <Box sx={{ width: '33%' }}> {/* Adjusted Box for FormControl */}
                                                <FormControl fullWidth>
                                                    <InputLabel id="time-unit-label">Override Workflow Time Unit</InputLabel>
                                                    <Select
                                                        labelId="or-time-unit-label"
                                                        id="time-unit-select"
                                                        value={customBasePeriodStepSizeUnit}
                                                        label="Time Unit"
                                                        onChange={(e) => setCustomBasePeriodStepSizeUnit(e.target.value)}
                                                    >
                                                        <MenuItem value="minutes">Minutes</MenuItem>
                                                        <MenuItem value="hours">Hours</MenuItem>
                                                        <MenuItem value="days">Days</MenuItem>
                                                        <MenuItem value="weeks">Weeks</MenuItem>
                                                    </Select>
                                                </FormControl>
                                            </Box>
                                            </Stack>
                                        </div>
                                    }
                                <Box sx={{ width: '50%', ml: 2 }}> {/* Adjusted Box for TextField */}
                                    <FormControlLabel
                                        control={<Switch checked={customBasePeriod} onChange={handleToggleChangePeriod} />}
                                        label={customBasePeriod ? 'Override Base Period' : 'Default Base Period' }
                                    />
                                </Box>
                                <Stack direction="row"  sx={{ ml: 2, mr: 2, mt: 4, mb: 2 }}>
                                    <Box sx={{ width: '50%', mr: 2 }}> {/* Adjusted Box for TextField */}
                                        <TextField
                                            type="number"
                                            label="Unix Start Time"
                                            variant="outlined"
                                            inputProps={{ min: 0 }}  // Set minimum value to 1
                                            value={unixStartTime}
                                            onChange={(e)=> setUnixStartTime(Number(e.target.value))}
                                            fullWidth
                                        />
                                    </Box>
                                    <Box sx={{ width: '50%', mr: 2 }}> {/* Adjusted Box for TextField */}
                                        <TextField
                                            type="datetime-local"
                                            variant="outlined"
                                            value={convertUnixToDateTimeLocal(unixStartTime)}
                                            onChange={handleDateTimeChange}
                                            fullWidth
                                        />
                                    </Box>
                                    <Box flexGrow={3} sx={{ mb: 0, mt: 1, mr: 2 }}>
                                        <Button fullWidth variant="outlined" onClick={() => unixStartTime > 0 ? setUnixStartTime(0) : setUnixStartTime(getCurrentUnixTimestamp())} >{ unixStartTime > 0 ? 'Reset' : 'Now'}</Button>
                                    </Box>
                                </Stack>
                                <Stack direction="row"  sx={{ ml: 2, mr: 0, mt: 4, mb: 2 }}>
                                    <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                        <Button fullWidth variant="outlined" onClick={() => 'Previous' === ti ? setUnixStartTime(unixStartTime-300) : setUnixStartTime(unixStartTime+300)} >{'Previous' === ti ? '-' : '+' } { ti} 5 Minutes</Button>
                                    </Box>
                                    <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                        <Button fullWidth variant="outlined" onClick={() => 'Previous' === ti ? setUnixStartTime(unixStartTime-3600) : setUnixStartTime(unixStartTime+3600)} >{'Previous' === ti ? '-' : '+' } { ti} 1 Hour</Button>
                                    </Box>
                                    <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                        <Button fullWidth variant="outlined" onClick={() =>  'Previous' === ti ? setUnixStartTime(unixStartTime-86400) : setUnixStartTime(unixStartTime+86400)} >{'Previous' === ti ? '-' : '+' }{ ti} 24 Hours</Button>
                                    </Box>
                                    <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                        <Button fullWidth variant="outlined" onClick={() => 'Previous' === ti ? setUnixStartTime(unixStartTime-604800) : setUnixStartTime(unixStartTime+604800)} >{'Previous' === ti ? '-' : '+' }{ ti} 7 Days</Button>
                                    </Box>
                                    <Box flexGrow={2} sx={{ mb: 2, mr: 2}}>
                                        <Button fullWidth variant="outlined" onClick={() => 'Previous' === ti ? setUnixStartTime(unixStartTime-2592000) : setUnixStartTime(unixStartTime+2592000)}>{'Previous' === ti ? '-' : '+' } {ti} 30 Days </Button>
                                    </Box>
                                </Stack>
                                <Box flexGrow={3} sx={{ mb: 2, ml: 1, mr: 1 }}>
                                    <Typography variant="h6" color="text.secondary">
                                        If you want it to start running immediately use 0. Otherwise set the time to when the workflow should start running, then select workflows from table below for scheduling.
                                    </Typography>
                                </Box>
                            </CardContent>
                        </Card>
                    </Container>
                    }
                    {requestStatus != '' && (
                        <Container sx={{  mt: 2}}>
                            <Typography variant="h6" color={requestStatusError}>
                                {requestStatus}
                            </Typography>
                        </Container>
                    )}
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                            <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                                <Tab label="Search" />
                                <Tab className="onboarding-card-highlight-all-search-indexer" label="Indexer"/>
                                <Tab className="onboarding-card-highlight-all-workflows" label="Workflows"/>
                                <Tab className="onboarding-card-highlight-all-workflows-runs" label="Runs"/>
                                {/*<Tab className="onboarding-card-highlight-all-workflows-actions" label="Actions"/>*/}
                            </Tabs>
                        </Box>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        { (selectedMainTab === 0) &&
                            <AiSearchAnalysis code={code} onChange={onChangeText} />
                        }
                        { (selectedMainTab === 1) &&
                            <SearchIndexersTable />
                        }
                        { (selectedMainTab === 2) &&
                            <div>
                                { selected && selected.length > 0  &&
                                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                        <Box sx={{ mb: 2 }}>
                                            <span>({selected.length} Selected Workflows)</span>
                                            <Button variant="outlined" color="secondary" onClick={(event) => handleWorkflowAction(event, 'start')} style={{marginLeft: '10px'}}>
                                                Start { selected.length === 1 ? 'Workflow' : 'Workflows' }
                                            </Button>
                                        </Box>
                                    </Container>
                                }
                                <WorkflowTable />
                            </div>
                        }
                        { (selectedMainTab === 3) &&
                            <div>

                                {requestRunsStatus != '' && (
                                    <Container sx={{  mt: 2}}>
                                        <Typography variant="h6" color={requestRunsStatusError}>
                                            {requestRunsStatus}
                                        </Typography>
                                    </Container>
                                )}
                                { selectedRuns && selectedRuns.length > 0  &&
                                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                        <Box sx={{ mb: 2 }}>
                                            <span>({selectedRuns.length} Selected Runs)</span>
                                            <Button variant="outlined" color="secondary" onClick={(event) => handleRunsActionRequest(event, 'stop')} style={{marginLeft: '10px'}}>
                                                Stop { selectedRuns.length === 1 ? 'Run' : 'Runs' }
                                            </Button>
                                        </Box>
                                    </Container>
                                }
                                <WorkflowAnalysisTable  />
                            </div>
                        }
                    </Container>
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}
type ValuePiece = Date | string | null;

type Value = ValuePiece | [ValuePiece, ValuePiece];
export default function AiWorkflowsDashboard() {
    return <AiWorkflowsDashboardContent />;
}
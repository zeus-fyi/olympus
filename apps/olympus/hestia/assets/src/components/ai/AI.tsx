import * as React from 'react';
import {useCallback, useState} from 'react';
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
    InputLabel,
    MenuItem,
    Select,
    Stack,
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
    setGroupFilter,
    setPlatformFilter,
    setSearchContent,
    setSearchResults,
    setUsernames,
    setWorkflowInstructions
} from "../../redux/ai/ai.reducer";
import {aiApiGateway} from "../../gateway/ai";
import {set} from 'date-fns';
import {TimeRange} from '@matiaslgonzalez/react-timeline-range-slider';

const mdTheme = createTheme();
const analysisStart = "====================================================================================ANALYSIS====================================================================================\n"
const analysisDone = "====================================================================================ANALYSIS-DONE===============================================================================\n"

function AiWorkflowsDashboardContent(props: any) {
    const [open, setOpen] = useState(true);
    const [loading, setIsLoading] = useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const searchKeywordsText = useSelector((state: RootState) => state.ai.searchContentText);
    const groupFilter = useSelector((state: RootState) => state.ai.groupFilter);
    const usernames = useSelector((state: RootState) => state.ai.usernames);
    const workflowInstructions = useSelector((state: RootState) => state.ai.workflowInstructions);
    const [code, setCode] = useState('');
    const searchResults = useSelector((state: RootState) => state.ai.searchResults);
    const platformFilter = useSelector((state: RootState) => state.ai.platformFilter);
    const [analyzeNext, setAnalyzeNext] = useState(false);
    const dispatch = useDispatch();
    const now = new Date();
    const getTodayAtSpecificHour = (hour: number = 12) =>
        set(now, { hours: hour, minutes: 0, seconds: 0, milliseconds: 0 });
    const [cycleCount, setCycleCount] = useState(0);
    const handleCycleCountChange = (event: any) => {
        setCycleCount(event.target.value);
    };
    const [aggregationCycleCount, setAggregationCycleCount] = useState(0);
    const handleAggregationCycleCountChange = (event: any) => {
        setAggregationCycleCount(event.target.value);
    };
    const [stepSize, setStepSize] = useState(5);
    const handleTimeStepChange = (event: any) => {
        setStepSize(event.target.value);
    };
    const [stepSizeUnit, setStepSizeUnit] = React.useState('minutes');
    const handleUpdateStepSizeUnit = (event: any) => {
        setStepSizeUnit(event.target.value);
    };
    const [analysisModel, setAnalysisModel] = React.useState('gpt-3.5-turbo');
    const handleUpdateAnalysisModel = (event: any) => {
        setAnalysisModel(event.target.value);
    };
    const [analysisModelMaxTokens, setAnalysisModelMaxTokens] = React.useState(0);
    const handleUpdateAnalysisModelMaxTokens = (event: any) => {
        setAnalysisModelMaxTokens(event.target.value);
    };
    const [analysisModelTokenOverflowStrategy, setAnalysisModelTokenOverflowStrategy] = React.useState('deduce');
    const handleUpdateAnalysisModelTokenOverflowStrategy = (event: any) => {
        setAnalysisModelTokenOverflowStrategy(event.target.value);
    };
    const [aggregationModel, setAggregationModel] = React.useState('gpt-4-1106-preview');
    const handleUpdateAggregationModel = (event: any) => {
        setAggregationModel(event.target.value);
    };
    const [workflowName, setWorkflowName] = React.useState('');
    const handleUpdateWorkflowName = (event: any) => {
        setWorkflowName(event.target.value);
    };

    const [aggregationModelTokenOverflowStrategy, setAggregationModelTokenOverflowStrategy] = React.useState('deduce');
    const handleUpdateAggregationModelTokenOverflowStrategy = (event: any) => {
        setAggregationModelTokenOverflowStrategy(event.target.value);
    };
    const [aggregationModelMaxTokens, setAggregationModelMaxTokens] = React.useState(0);
    const handleUpdateAggregationModelMaxTokens = (event: any) => {
        setAggregationModelMaxTokens(event.target.value);
    };
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
    const handleUpdateSearchKeywords = (value: string) => {
        dispatch(setSearchContent(value));
    };
    const handleUpdateGroupFilter = (value: string) => {
        dispatch(setGroupFilter(value));
    };
    const handleUpdatePlatformFilter = (value: string) => {
        dispatch(setPlatformFilter(value));
    };
    const handleUpdateSearchUsernames =(value: string) => {
        dispatch(setUsernames(value));
    };
    const handleUpdateWorkflowInstructions =(value: string) => {
        dispatch(setWorkflowInstructions(value));
    };

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    const handleSearchRequest = async (timeRange: '1 hour'| '24 hours' | '7 days'| '30 days' | 'window' | 'all') => {
        try {
            setIsLoading(true)
            console.log(searchInterval, 'sdfs')
            const response = await aiApiGateway.searchRequest({
                'searchContentText': searchKeywordsText,
                'groupFilter': groupFilter,
                'platforms': platformFilter,
                'usernames': usernames,
                'workflowInstructions': workflowInstructions,
                'searchInterval': searchInterval,
                'timeRange': timeRange,
            });
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSearchResults(data));
                setCode(data)
            } else {
                console.log('Failed to search', response);
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    }
    const handleSearchAnalyzeRequest = async (timeRange: '1 hour'| '24 hours' | '7 days'| '30 days' | 'window' | 'all') => {
        try {
            setIsLoading(true)
            const response = await aiApiGateway.analyzeSearchRequest({
                'searchContentText': searchKeywordsText,
                'groupFilter': groupFilter,
                'platforms': platformFilter,
                'usernames': usernames,
                'workflowInstructions': workflowInstructions,
                'searchInterval': searchInterval,
                'cycleCount': cycleCount,
                'aggregationCycleCount': aggregationCycleCount,
                'stepSize': stepSize,
                'stepSizeUnit': stepSizeUnit,
                'timeRange': timeRange,
                'workflowName': workflowName,
                'analysisModel': analysisModel,
                'analysisModelMaxTokens': analysisModelMaxTokens,
                'analysisModelTokenOverflowStrategy': analysisModelTokenOverflowStrategy,
                'aggregationModel': aggregationModel,
                'aggregationModelMaxTokens': aggregationModelMaxTokens,
                'aggregationModelTokenOverflowStrategy': aggregationModelTokenOverflowStrategy,
            });
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSearchResults(data));
                setCode(analysisStart + data+ analysisDone + code)
            } else {
                console.log('Failed to search', response);
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    }
    if (loading) {
        return <div>Loading...</div>;
    }

    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    };
    const onChangeText = (textInput: string) => {
        setCode(textInput);
    };

    const formatTick2 = (ms: number) => {
        return new Date(ms).toLocaleTimeString([], { hour: 'numeric', hour12: true });
    };

    const ti = analyzeNext ? 'Next' : 'Previous';
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
                            LLM Workflow Engine
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
                        <Stack direction="row" spacing={2}>
                            <Card sx={{ minWidth: 100, maxWidth: 600 }}>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Search Augmented LLM Workflow Engine
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Currently limited public functionality. This will allow you to search across many platforms including our own for data, or give the
                                        AI enough context to build its own workflows that can map-reduce analyze, run devops tasks, or even build apps.
                                        In the meantime you can email ai@zeus.fyi and it'll summarize your email, and suggest responses.
                                        When adding many values to a field use comma delimited entries.
                                    </Typography>
                                </CardContent>
                                <CardContent>
                                    <Stack direction="column" >
                                        <Box flexGrow={1} sx={{ mb: 2 }}>
                                            <TextField
                                                fullWidth
                                                id="platforms-input"
                                                label="Platforms"
                                                variant="outlined"
                                                value={platformFilter}
                                                onChange={(e) => handleUpdatePlatformFilter(e.target.value)}
                                            />
                                        </Box><Box flexGrow={1} sx={{ mb: 2 }}>
                                            <TextField
                                                fullWidth
                                                id="group-input"
                                                label="Group"
                                                variant="outlined"
                                                value={groupFilter}
                                                onChange={(e) => handleUpdateGroupFilter(e.target.value)}
                                            />
                                        </Box>
                                        <Box flexGrow={1} sx={{ mb: 2 }}>
                                            <TextField
                                                fullWidth
                                                id="usernames-input"
                                                label="Usernames"
                                                variant="outlined"
                                                value={usernames}
                                                onChange={(e) => handleUpdateSearchUsernames(e.target.value)}
                                            />
                                        </Box>
                                        <Box flexGrow={1} sx={{ mb: 2 }}>
                                            <TextField
                                                fullWidth
                                                id="keywords-input"
                                                label="Content"
                                                variant="outlined"
                                                value={searchKeywordsText}
                                                onChange={(e) => handleUpdateSearchKeywords(e.target.value)}
                                            />
                                        </Box>
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
                                    </Stack>
                                </CardContent>
                            </Card>
                            <Card sx={{ minWidth: 500, maxWidth: 900 }}>
                                <CardContent>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Workflow Generation
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        This allows you to write natural language instructions to chain to your search queries. Add a name
                                        for your workflow, and then write instructions for the AI to follow, and it will save the workflow for you.
                                    </Typography>
                                </CardContent>
                                <CardContent>
                                    <Box sx={{ width: '100%', mb: 2, mt: -2 }}>
                                        <TextField
                                            label={`Workflow Name`}
                                            variant="outlined"
                                            value={workflowName}
                                            onChange={handleUpdateWorkflowName}
                                            fullWidth
                                        />
                                    </Box>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Analysis Instructions
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Token overflow strategy will determine how the AI will handle requests that are projected to exceed the maximum token length for the model you select, or has returned a result with that error.
                                        Deduce will chunk your analysis into smaller pieces and aggregate them into a final analysis result. Truncate will simply truncate the request
                                        to the maximum token length it can support.
                                    </Typography>
                                    <Stack direction="row" >
                                        <Box flexGrow={3} sx={{ mb: 4, mt: 4 }}>
                                            <FormControl fullWidth>
                                                <InputLabel id="analysis-model-label">Analysis Model</InputLabel>
                                                <Select
                                                    labelId="analysis-model-label"
                                                    id="analysis-model-select"
                                                    value={analysisModel}
                                                    label="Analysis Model"
                                                    onChange={handleUpdateAnalysisModel}
                                                >
                                                    <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                    <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <Box flexGrow={2} sx={{ mb: 4, mt: 4, ml:2 }}>
                                            <FormControl fullWidth>
                                                <InputLabel id="token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                <Select
                                                    labelId="token-overflow-analysis-label"
                                                    id="analysis-overflow-analysis-select"
                                                    value={analysisModelTokenOverflowStrategy}
                                                    label="Token Overflow Strategy"
                                                    onChange={handleUpdateAnalysisModelTokenOverflowStrategy}
                                                >
                                                    <MenuItem value="deduce">deduce</MenuItem>
                                                    <MenuItem value="truncate">truncate</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <Box flexGrow={1} sx={{ mb: 4, mt: 4, ml:2 }}>
                                            <TextField
                                                type="number"
                                                label="Cycle Count"
                                                variant="outlined"
                                                value={cycleCount}
                                                inputProps={{ min: 0 }}  // Set minimum value to 0
                                                onChange={handleCycleCountChange}
                                                fullWidth
                                            />
                                        </Box>
                                    </Stack>
                                    <Box  sx={{ mb: 2, mt: -2 }}>
                                        <TextareaAutosize
                                            minRows={18}
                                            value={workflowInstructions}
                                            onChange={(e) => handleUpdateWorkflowInstructions(e.target.value)}
                                            style={{ resize: "both", width: "100%" }}
                                        />
                                    </Box>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Aggregation Instructions
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Use this to tell the AI how to aggregate the results of your analysis chunks into a rolling aggregation window. This is useful for
                                        allowing you to create higher level analysis on top of your search results that isn't possible or desired in a single cycle.
                                        Aggregation cycle count sets how many analysis cycles per aggregation cycle. For example if you set 10 total analysis cycles, and 2 aggregation cycles,
                                        it will aggregate the results from each 5 analysis cycles into a single aggregation cycle,
                                        and use that as the input for the next aggregation cycle. It will run a final aggregation cycle after the last analysis cycle if
                                        the total analysis cycle doesn't divide evenly into the aggregation count.
                                    </Typography>
                                    <Stack direction="row" >
                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                            <FormControl fullWidth>
                                                <InputLabel id="aggregation-model-label">Aggregation Model</InputLabel>
                                                <Select
                                                    labelId="aggregation-model-label"
                                                    id="aggregation-model-select"
                                                    value={aggregationModel}
                                                    label="Aggregation Model"
                                                    onChange={handleUpdateAggregationModel}
                                                >
                                                    <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                    <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                            <FormControl fullWidth>
                                                <InputLabel id="aggregation-token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                <Select
                                                    labelId="aggregation-token-overflow-analysis-label"
                                                    id="aggregation-token-overflow-analysis-select"
                                                    value={aggregationModelTokenOverflowStrategy}
                                                    label="Token Overflow Strategy"
                                                    onChange={handleUpdateAggregationModelTokenOverflowStrategy}
                                                >
                                                    <MenuItem value="deduce">deduce</MenuItem>
                                                    <MenuItem value="truncate">truncate</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                            <TextField
                                                type="number"
                                                label="Aggregation Cycle Count"
                                                variant="outlined"
                                                value={aggregationCycleCount}
                                                inputProps={{ min: 0, max: cycleCount }}  // Set minimum value to 0
                                                onChange={handleAggregationCycleCountChange}
                                                fullWidth
                                            />
                                        </Box>
                                    </Stack>
                                    <Box  sx={{ mb: 2, mt: 2 }}>
                                        <TextareaAutosize
                                            minRows={18}
                                            value={workflowInstructions}
                                            onChange={(e) => handleUpdateWorkflowInstructions(e.target.value)}
                                            style={{ resize: "both", width: "100%" }}
                                        />
                                    </Box>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Time Intervals
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        You can run an analysis on demand or use this to define an analysis chunk interval as part of an aggregate analysis.
                                    </Typography>
                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                        <Box sx={{ width: '33%' }}> {/* Adjusted Box for TextField */}
                                            <TextField
                                                type="number"
                                                label="Time Step Size"
                                                variant="outlined"
                                                inputProps={{ min: 1 }}  // Set minimum value to 1
                                                value={stepSize}
                                                onChange={handleTimeStepChange}
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
                                                    onChange={handleUpdateStepSizeUnit}
                                                >
                                                    <MenuItem value="seconds">Seconds</MenuItem>
                                                    <MenuItem value="minutes">Minutes</MenuItem>
                                                    <MenuItem value="hours">Hours</MenuItem>
                                                    <MenuItem value="days">Days</MenuItem>
                                                    <MenuItem value="weeks">Weeks</MenuItem>
                                                </Select>
                                            </FormControl>
                                        </Box>
                                        <Box sx={{ width: '33%', mb: 4 }}> {/* New Box for Total Time TextField */}
                                            <TextField
                                                label={`Total Time (${stepSizeUnit})`} // Label now reflects the selected unit
                                                variant="outlined"
                                                value={stepSize* cycleCount}
                                                InputProps={{
                                                    readOnly: true,
                                                }}
                                                fullWidth
                                            />
                                        </Box>
                                    </Stack>
                                    <Typography gutterBottom variant="h5" component="div">
                                        Workflow Token Usage Limits
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Use this to define how many tokens you want to use for your analysis. This will allow you to control how much you spend on your analysis.
                                        Set to 0 for unlimited, otherwise it will stop the analysis when it reaches the limit per model.
                                    </Typography>
                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                        <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                            <TextField
                                                label={`Analysis Model`}
                                                variant="outlined"
                                                value={analysisModel}
                                                InputProps={{
                                                    readOnly: true,
                                                }}
                                                fullWidth
                                            />
                                        </Box>
                                        <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                            <TextField
                                                type="number"
                                                label={`Max Tokens Analysis Model`}
                                                variant="outlined"
                                                value={analysisModelMaxTokens}
                                                inputProps={{ min: 0 }}
                                                onChange={handleUpdateAnalysisModelMaxTokens}
                                                fullWidth
                                            />
                                        </Box>
                                    </Stack>
                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                        <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                            <TextField
                                                label={`Aggregation Model`}
                                                variant="outlined"
                                                value={aggregationModel}
                                                InputProps={{
                                                    readOnly: true,
                                                }}
                                                fullWidth
                                            />
                                        </Box>
                                        <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                            <TextField
                                                type="number"
                                                label={`Max Aggregation Token Usage`}
                                                variant="outlined"
                                                value={aggregationModelMaxTokens}
                                                onChange={handleUpdateAggregationModelMaxTokens}
                                                inputProps={{ min: 0 }}
                                                fullWidth
                                            />
                                        </Box>
                                    </Stack>
                                    <Box flexGrow={1} sx={{ mb: 2 }}>
                                        <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('all')} >Save Workflow</Button>
                                    </Box>
                                    {/*<Typography gutterBottom variant="h5" component="div">*/}
                                    {/*    Workflow Generations*/}
                                    {/*</Typography>*/}
                                    {/*<Typography variant="body2" color="text.secondary">*/}
                                    {/*    Use Start Working Analysis to generate a workflow that will run the analysis on the time intervals you've defined. It will*/}
                                    {/*    process the data that gets generated from your search query, and then aggregate the results into a rolling window.*/}
                                    {/*</Typography>*/}
                                    {/*<Box flexGrow={1} sx={{ mt: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('all')} >Start Working Analysis</Button>*/}
                                    {/*</Box>*/}
                                    {/*<Box flexGrow={1} sx={{ mt: 2 }}>*/}
                                    {/*    <Typography variant="body2" color="text.secondary">*/}
                                    {/*        Use these buttons to search previous time intervals relative to the current time.*/}
                                    {/*    </Typography>*/}
                                    {/*</Box>*/}
                                    {/*<FormControlLabel*/}
                                    {/*    control={<Switch checked={analyzeNext} onChange={handleToggleChange} />}*/}
                                    {/*    label={analyzeNext ? 'Analyze Next' : 'Analyze Previous'}*/}
                                    {/*/>*/}
                                    {/*<Box flexGrow={1} sx={{ mb: 2, mt: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('1 hour')} >Analyze {ti} 1 Hour</Button>*/}
                                    {/*</Box>*/}
                                    {/*<Box flexGrow={1} sx={{ mb: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('24 hours')} >Analyze {ti} 24 Hours</Button>*/}
                                    {/*</Box>*/}
                                    {/*<Box flexGrow={1} sx={{ mb: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('7 days')} >Analyze {ti} 7 Days</Button>*/}
                                    {/*</Box>*/}
                                    {/*<Box flexGrow={1} sx={{ mb: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('30 days')} >Analyze {ti} 30 Days </Button>*/}
                                    {/*</Box>*/}
                                    {/*<Box flexGrow={1} sx={{ mb: 2 }}>*/}
                                    {/*    <Button fullWidth variant="contained" onClick={() => handleSearchAnalyzeRequest('all')} >Analyze All {ti} Records</Button>*/}
                                    {/*</Box>*/}
                                </CardContent>
                            </Card>
                        </Stack>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                            <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                                <Tab label="Search" />
                                <Tab className="onboarding-card-highlight-all-workflows" label="Workflows"  />
                            </Tabs>
                        </Box>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        { (selectedMainTab === 0) &&
                            <AiSearchAnalysis code={code} onChange={onChangeText} />
                        }
                        { (selectedMainTab === 1) &&
                            <WorkflowTable loading={loading}/>
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
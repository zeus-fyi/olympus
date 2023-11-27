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
    const searchContentText = useSelector((state: RootState) => state.ai.searchContentText);
    const groupFilter = useSelector((state: RootState) => state.ai.groupFilter);
    const usernames = useSelector((state: RootState) => state.ai.usernames);
    const workflowInstructions = useSelector((state: RootState) => state.ai.workflowInstructions);
    const [code, setCode] = useState('');
    const searchResults = useSelector((state: RootState) => state.ai.searchResults);
    const platformFilter = useSelector((state: RootState) => state.ai.platformFilter);
    const dispatch = useDispatch();
    const now = new Date();
    const getTodayAtSpecificHour = (hour: number = 12) =>
        set(now, { hours: hour, minutes: 0, seconds: 0, milliseconds: 0 });
    const [cycleCount, setCycleCount] = useState(0);
    const handleCycleCountChange = (event: any) => {
        setCycleCount(event.target.value);
    };
    const [stepSize, setStepSize] = useState(5);
    const handleTimeStepChange = (event: any) => {
        setStepSize(event.target.value);
    };
    const [stepSizeUnit, setStepSizeUnit] = React.useState('minutes');
    const handleUpdateStepSizeUnit = (event: any) => {
        setStepSizeUnit(event.target.value);
    };
    const [searchInterval, setSearchInterval] = useState<[Date, Date]>([getTodayAtSpecificHour(0), getTodayAtSpecificHour(24)]);
    const onTimeRangeChange = useCallback((interval: [Date, Date]) => {
        setSearchInterval(interval);
    }, []);

    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();

    const handleUpdateSearchContent = (value: string) => {
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

    const handleSearchRequest = async (timeRange: '24 hours' | '7 days'| '30 days' | 'window' | 'all') => {
        try {
            setIsLoading(true)
            console.log(searchInterval, 'sdfs')
            const response = await aiApiGateway.searchRequest({
                'searchContentText': searchContentText,
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

    const handleSearchAnalyzeRequest = async () => {
        try {
            setIsLoading(true)
            const response = await aiApiGateway.analyzeSearchRequest({
                'searchContentText': searchContentText,
                'groupFilter': groupFilter,
                'platforms': platformFilter,
                'usernames': usernames,
                'workflowInstructions': workflowInstructions,
                'searchInterval': searchInterval,
                'cycleCount': cycleCount,
                'stepSize': stepSize,
                'stepSizeUnit': stepSizeUnit,
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
                                                id="content-input"
                                                label="Content"
                                                variant="outlined"
                                                value={searchContentText}
                                                onChange={(e) => handleUpdateSearchContent(e.target.value)}
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
                                        <Box flexGrow={1} sx={{ mb: 2, mt: 2 }}>
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
                                        Workflow Instructions
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        This allows you to write natural language instructions to chain to your search queries.
                                    </Typography>
                                </CardContent>
                                <CardContent>
                                    <Box  sx={{ mb: 2 }}>
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
                                        <Box sx={{ width: '50%' }}> {/* Adjusted Box for TextField */}
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
                                        <Box sx={{ width: '30%' }}> {/* Adjusted Box for FormControl */}
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
                                        <Box sx={{ width: '20%'}}> {/* Adjusted Box for Cycle Count TextField */}
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
                                    <Box sx={{ width: '100%', mb: 4 }}> {/* New Box for Total Time TextField */}
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
                                    <Button fullWidth variant="contained" onClick={handleSearchAnalyzeRequest} >Start Working Analysis</Button>
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